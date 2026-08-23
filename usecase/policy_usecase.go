package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
	"github.com/janghanul090801/pico-backend/external/ai"
	"github.com/janghanul090801/pico-backend/external/kakaomap"
	"github.com/janghanul090801/pico-backend/external/volunteer"
	"github.com/janghanul090801/pico-backend/external/welfare"
	"github.com/janghanul090801/pico-backend/external/youth"
	"github.com/janghanul090801/pico-backend/internal/collections"
)

type policyUseCase struct {
	policyRepository      domain.PolicyRepository
	policyMatchRepository domain.PolicyMatchRepository
	userRepository        domain.UserRepository
	youthClient           *youth.Client
	welfareClient         *welfare.Client
	volunteerClient       *volunteer.Client
	kakaoMapClient        *kakaomap.Client
	aiClient              *ai.Client
	matcher               *PolicyMatcher
	contextTimeout        time.Duration
}

func NewPolicyUseCase(
	policyRepository domain.PolicyRepository,
	policyMatchRepository domain.PolicyMatchRepository,
	userRepository domain.UserRepository,
	youthClient *youth.Client,
	welfareClient *welfare.Client,
	volunteerClient *volunteer.Client,
	kakaoMapClient *kakaomap.Client,
	aiClient *ai.Client,
	matcher *PolicyMatcher,
	contextTimeout time.Duration,
) domain.PolicyUseCase {
	return &policyUseCase{
		policyRepository:      policyRepository,
		policyMatchRepository: policyMatchRepository,
		userRepository:        userRepository,
		youthClient:           youthClient,
		welfareClient:         welfareClient,
		volunteerClient:       volunteerClient,
		kakaoMapClient:        kakaoMapClient,
		aiClient:              aiClient,
		matcher:               matcher,
		contextTimeout:        contextTimeout,
	}
}

func (u *policyUseCase) GetByID(c context.Context, id *domain.ID) (*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	p, err := u.policyRepository.FindByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewNotFoundError(err)
		}
		return nil, domain.NewInternalServerError(err)
	}

	return p, nil
}

func (u *policyUseCase) ListByRegionCodeAndActive(c context.Context, regionCode int) ([]*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	policies, err := u.policyRepository.FindAllByRegionCodeAndActive(ctx, regionCode)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}
	return policies, nil
}

func (u *policyUseCase) ListMatchesByUserID(c context.Context, userID *domain.ID) ([]*domain.PolicyMatch, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	matches, err := u.policyMatchRepository.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return matches, nil
}

func (u *policyUseCase) GetFromApi(c context.Context) ([]*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout*50)
	defer cancel()

	var fetchedPolicies []*domain.Policy

	youthPolicies, err := u.youthClient.Fetch(ctx)
	if err != nil {
		log.Printf("[PolicyUseCase] youch fetch warning: %v", err)
	}
	welfarePolicies, err := u.welfareClient.Fetch(ctx)
	if err != nil {
		log.Printf("[PolicyUseCase] welfare fetch warning: %v", err)
	}
	volunteerPolicies, err := u.volunteerClient.Fetch(ctx)
	if err != nil {
		log.Printf("[PolicyUseCase] volunteer fetch warning: %v", err)
	}
	// maternity, _ := u.maternityClient.Fetch(ctx)

	policies := append(youthPolicies, welfarePolicies...)
	policies = append(policies, volunteerPolicies...)
	// policies = append(policies, maternity...)

	log.Printf("[PolicyUseCase] Total fetched raw policies: %d", len(fetchedPolicies))

	seoulPolicies := collections.Filter(policies, func(p *domain.Policy) bool {
		return p.IsSeoulOrNational()
	})

	log.Printf("[PolicyUseCase] Filtered policies for Seoul target: %d", len(seoulPolicies))

	var savedPolicies []*domain.Policy
	for _, p := range seoulPolicies {
		lat, long, err := u.kakaoMapClient.Geocode(ctx, p.Address)
		if err != nil {
			log.Printf("[PolicyUseCase] policy id: %s, kakaoMapClient.Geocode error: %v", p.ID, err)
		}
		p.Latitude = lat
		p.Longitude = long
		interests, _ := u.getInterests(c, p)
		p.Condition.Interests = interests
		saved, err := u.policyRepository.Create(ctx, p)
		if err != nil {
			log.Printf("[PolicyUseCase] Failed to save policy '%s': %v", p.Title, err)
			continue
		}
		savedPolicies = append(savedPolicies, saved)
	}

	return savedPolicies, nil
}

func (u *policyUseCase) List(c context.Context) ([]*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.policyRepository.FindAll(ctx)
}

func (u *policyUseCase) SavePolicyMatches(c context.Context, user *domain.User, policies []*domain.Policy) ([]*domain.PolicyMatch, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout*100)
	defer cancel()

	var match *domain.PolicyMatch
	var err error
	matches := make([]*domain.PolicyMatch, len(policies))
	for i, policy := range policies {
		status := u.matcher.Match(user, policy)
		match = &domain.PolicyMatch{
			PolicyID: policy.ID,
			UserID:   user.ID,
			Status:   status,
		}
		if status == domain.MatchUNCERTAIN {
			match.Probability, match.Comment, _ = u.getProbabilityAndComment(c, user, policy)
		}
		matches[i], err = u.policyMatchRepository.Create(ctx, match)
		if err != nil {
			return nil, domain.NewInternalServerError(err)
		}
	}

	return matches, nil
}

func (u *policyUseCase) UpdatePolicyMatchesForUser(c context.Context, userID *domain.ID) ([]*domain.PolicyMatch, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout*100)
	defer cancel()

	usr, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	policies, err := u.policyRepository.FindAll(ctx)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	var updatedMatches []*domain.PolicyMatch
	for _, policy := range policies {
		status := u.matcher.Match(usr, policy)
		var prob *int
		var comm *string
		if status == domain.MatchUNCERTAIN {
			prob, comm, _ = u.getProbabilityAndComment(c, usr, policy)
		}

		existing, err := u.policyMatchRepository.FindByUserIDAndPolicyID(ctx, userID, &policy.ID)
		if err != nil || existing == nil {
			newMatch := &domain.PolicyMatch{
				PolicyID:    policy.ID,
				UserID:      usr.ID,
				Status:      status,
				Probability: prob,
				Comment:     comm,
			}
			created, createErr := u.policyMatchRepository.Create(ctx, newMatch)
			if createErr == nil {
				updatedMatches = append(updatedMatches, created)
			}
		} else {
			updated, updateErr := u.policyMatchRepository.Update(ctx, &existing.ID, status, prob, comm)
			if updateErr == nil {
				updatedMatches = append(updatedMatches, updated)
			}
		}
	}

	return updatedMatches, nil
}

func (u *policyUseCase) getProbabilityAndComment(c context.Context, user *domain.User, policy *domain.Policy) (*int, *string, error) {
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	userInfo := string(data)
	data, err = json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	policyInfo := string(data)
	res, err := u.aiClient.Chat(c, fmt.Sprintf(domain.MatchPrompt, userInfo, policyInfo))
	if err != nil {
		return nil, nil, err
	}

	resString := strings.Split(res, ",")

	probability, err := strconv.Atoi(resString[0])
	if err != nil {
		return nil, nil, err
	}

	comment := resString[1]
	return &probability, &comment, nil
}

func (u *policyUseCase) getInterests(c context.Context, policy *domain.Policy) ([]domain.Interest, error) {
	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return nil, err
	}

	policyInfo := string(data)
	res, err := u.aiClient.Chat(c, fmt.Sprintf(domain.InterestPrompt, policyInfo))
	if err != nil {
		fmt.Printf("AI API ERROR: %v\n", err)
		return nil, err
	}
	return collections.Filter(collections.Map(strings.Split(res, ","), func(s string) domain.Interest {
		return domain.Interest(s)
	}), func(i domain.Interest) bool {
		switch i {
		case domain.InterestEmployment:
			return true
		case domain.InterestHousing:
			return true
		case domain.InterestEducation:
			return true
		case domain.InterestWelfare:
			return true
		case domain.InterestPregnancy:
			return true
		case domain.InterestCulture:
			return true
		case domain.InterestEnvironment:
			return true
		case domain.InterestParticipation:
			return true
		default:
			return true
		}
	}), nil
}

//// FetchMaternityPolicies fetches maternity & childcare support status (도/시 출산장려/양육비 지원현황)
//func (u *policyUseCase) FetchMaternityPolicies(ctx context.Context) ([]*domain.Policy, error) {
//	apiKey := config.E.MaternityApiKey
//
//	reqURL := fmt.Sprintf("http://apis.data.go.kr/1741000/MaternityChildcareService/getMaternityChildcareList?serviceKey=%s&ctpvNm=%s&pageNo=1&numOfRows=20", url.QueryEscape(apiKey), url.QueryEscape("서울특별시"))
//
//	body, err := u.doGetRequest(ctx, reqURL)
//	if err != nil {
//		return nil, fmt.Errorf("http request failed: %w", err)
//	}
//
//	var items []domain.MaternityItem
//
//	var xmlResp domain.MaternityXMLResponse
//	if xmlErr := xml.Unmarshal(body, &xmlResp); xmlErr == nil && len(xmlResp.Body.Items.ItemList) > 0 {
//		items = xmlResp.Body.Items.ItemList
//	} else {
//		var jsonResp domain.MaternityJSONResponse
//		if jsonErr := json.Unmarshal(body, &jsonResp); jsonErr == nil && len(jsonResp.Response.Body.Items.Item) > 0 {
//			items = jsonResp.Response.Body.Items.Item
//		}
//	}
//
//	var result []*domain.Policy
//	now := time.Now()
//
//	for _, item := range items {
//		title := strings.TrimSpace(item.Bznm)
//		if title == "" {
//			title = strings.TrimSpace(item.ServNm)
//		}
//		if title == "" {
//			continue
//		}
//
//		desc := strings.TrimSpace(item.DetlCn)
//		if desc == "" {
//			desc = strings.TrimSpace(item.ServDgst)
//		}
//		if desc == "" {
//			desc = title
//		}
//
//		addr := strings.TrimSpace(fmt.Sprintf("%s %s", item.CtpvNm, item.SggNm))
//		regionCode := parseRegionCode("", addr)
//
//		sDate := parseYYYYMMDD(item.Bgnde)
//		if sDate.IsZero() {
//			sDate = now
//		}
//		eDate := parseYYYYMMDD(item.Endde)
//		if eDate.IsZero() {
//			eDate = now.AddDate(1, 0, 0)
//		}
//
//		result = append(result, &domain.Policy{
//			Title:       title,
//			Description: desc,
//			RegionCode:  regionCode,
//			StartDate:   sDate,
//			EndDate:     eDate,
//			Address:     addr,
//			Latitude:    0.0,
//			Longitude:   0.0,
//		})
//	}
//
//	return result, nil
//}

// ReclassifyAllPolicies 는 이미 저장된 정책들 중 카테고리(interests) 분류가 안 된 것들을
// 다시 돌면서 AI로 재분류하는 일회성 관리자용 기능이다.
// GetFromApi의 자동 분류는 "새로 가져온 정책"에만 적용되기 때문에,
// 예전에 분류 없이 저장된 기존 데이터는 이 함수로 별도 재실행해야 한다.
func (u *policyUseCase) ReclassifyAllPolicies(c context.Context) (int, error) {
	policies, err := u.policyRepository.FindAll(c)
	if err != nil {
		return 0, domain.NewInternalServerError(err)
	}

	updated := 0
	for i, p := range policies {
		if len(p.Condition.Interests) > 0 {
			continue // 이미 분류되어 있으면 건너뜀
		}

		// 144개를 연달아 바로 쏘면 AI API의 rate limit(429)에 걸림.
		// 요청 사이에 약간의 텀을 줘서 속도를 늦춤.
		if i > 0 {
			time.Sleep(1500 * time.Millisecond)
		}

		ctx, cancel := context.WithTimeout(c, u.contextTimeout*10)
		interests, err := u.getInterests(ctx, p)
		cancel()
		if err != nil {
			log.Printf("[ReclassifyAllPolicies] policy id: %s, getInterests error: %v", p.ID, err)
			continue
		}
		if len(interests) == 0 {
			continue
		}

		newCondition := p.Condition
		newCondition.Interests = interests

		updateCtx, updateCancel := context.WithTimeout(c, u.contextTimeout)
		err = u.policyRepository.UpdateCondition(updateCtx, &p.ID, newCondition)
		updateCancel()
		if err != nil {
			log.Printf("[ReclassifyAllPolicies] policy id: %s, UpdateCondition error: %v", p.ID, err)
			continue
		}

		updated++
	}

	return updated, nil
}
