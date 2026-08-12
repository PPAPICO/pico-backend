package usecase

import (
	"context"
	"log"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
	"github.com/janghanul090801/pico-backend/external/kakaomap"
	"github.com/janghanul090801/pico-backend/external/volunteer"
	"github.com/janghanul090801/pico-backend/external/welfare"
	"github.com/janghanul090801/pico-backend/external/youth"
	"github.com/janghanul090801/pico-backend/internal/collections"
)

type policyUseCase struct {
	policyRepository      domain.PolicyRepository
	policyMatchRepository domain.PolicyMatchRepository
	youthClient           *youth.Client
	welfareClient         *welfare.Client
	volunteerClient       *volunteer.Client
	kakaoMapClient        *kakaomap.Client
	matcher               *PolicyMatcher
	contextTimeout        time.Duration
}

func NewPolicyUseCase(policyRepository domain.PolicyRepository, policyMatchRepository domain.PolicyMatchRepository, youthClient *youth.Client, welfareClient *welfare.Client, volunteerClient *volunteer.Client, kakaoMapClient *kakaomap.Client, matcher *PolicyMatcher, contextTimeout time.Duration) domain.PolicyUseCase {
	return &policyUseCase{
		policyRepository:      policyRepository,
		policyMatchRepository: policyMatchRepository,
		youthClient:           youthClient,
		welfareClient:         welfareClient,
		volunteerClient:       volunteerClient,
		kakaoMapClient:        kakaoMapClient,
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

// GetFromApi fetches policies from external APIs, filters for Seoul/National targets, and saves them to DB
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
