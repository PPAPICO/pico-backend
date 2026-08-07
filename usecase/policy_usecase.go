package usecase

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"golang.org/x/sync/errgroup"
)

type policyUseCase struct {
	policyRepository      domain.PolicyRepository
	policyMatchRepository domain.PolicyMatchRepository
	contextTimeout        time.Duration
	httpClient            *http.Client
}

func NewPolicyUseCase(policyRepository domain.PolicyRepository, policyMatchRepository domain.PolicyMatchRepository, contextTimeout time.Duration) domain.PolicyUseCase {
	return &policyUseCase{
		policyRepository:      policyRepository,
		policyMatchRepository: policyMatchRepository,
		contextTimeout:        contextTimeout,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (u *policyUseCase) GetByID(c context.Context, id *domain.ID) (*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.policyRepository.FindByID(ctx, id)
}

func (u *policyUseCase) ListRegionCodeAndActive(c context.Context, regionCode int) ([]*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.policyRepository.FindAllByRegionCodeAndActive(ctx, regionCode)
}

func (u *policyUseCase) GetMatchesByUserID(c context.Context, userID *domain.ID) ([]*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	matches, err := u.policyMatchRepository.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var policies []*domain.Policy
	for _, match := range matches {
		p, err := u.policyRepository.FindByID(ctx, &match.PolicyID)
		if err == nil && p != nil {
			policies = append(policies, p)
		}
	}
	return policies, nil
}

// GetFromApi fetches policies from external APIs, filters for Seoul/National targets, and saves them to DB
func (u *policyUseCase) GetFromApi(c context.Context) ([]*domain.Policy, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout*5)
	defer cancel()

	var fetchedPolicies []*domain.Policy

	// 1. 온통청년 청년정책 API (Youth Center)
	youthPolicies, err := u.FetchYouthPolicies(ctx)
	if err != nil {
		log.Printf("[PolicyUseCase] FetchYouthPolicies warning: %v", err)
	} else {
		fetchedPolicies = append(fetchedPolicies, youthPolicies...)
	}

	// 2. 한국사회보장정보원 복지서비스정보 (중앙부처복지서비스)
	welfarePolicies, err := u.FetchWelfarePolicies(ctx)
	if err != nil {
		log.Printf("[PolicyUseCase] FetchWelfarePolicies warning: %v", err)
	} else {
		fetchedPolicies = append(fetchedPolicies, welfarePolicies...)
	}

	// 3. 행정안전부 봉사참여정보서비스 (지역 행사/봉사활동)
	volunteerPolicies, err := u.FetchVolunteerEvents(ctx)
	if err != nil {
		log.Printf("[PolicyUseCase] FetchVolunteerEvents warning: %v", err)
	} else {
		fetchedPolicies = append(fetchedPolicies, volunteerPolicies...)
	}

	// 4. 도/시 출산장려/양육비 지원현황 (임산부/출산)
	maternityPolicies, err := u.FetchMaternityPolicies(ctx)
	if err != nil {
		log.Printf("[PolicyUseCase] FetchMaternityPolicies warning: %v", err)
	} else {
		fetchedPolicies = append(fetchedPolicies, maternityPolicies...)
	}

	log.Printf("[PolicyUseCase] Total fetched raw policies: %d", len(fetchedPolicies))

	// Strict filter for Seoul region (RegionCode == 11000 or National Policy RegionCode == 0)
	var seoulPolicies []*domain.Policy
	for _, p := range fetchedPolicies {
		if isSeoulOrNational(p) {
			seoulPolicies = append(seoulPolicies, p)
		}
	}

	log.Printf("[PolicyUseCase] Filtered policies for Seoul target: %d", len(seoulPolicies))

	// Save filtered policies to DB repository
	var savedPolicies []*domain.Policy
	for _, p := range seoulPolicies {
		saved, err := u.policyRepository.Create(ctx, p)
		if err != nil {
			log.Printf("[PolicyUseCase] Failed to save policy '%s': %v", p.Title, err)
			continue
		}
		savedPolicies = append(savedPolicies, saved)
	}

	return savedPolicies, nil
}

// FetchYouthPolicies fetches policies from 온통청년 청년정책 API
func (u *policyUseCase) FetchYouthPolicies(ctx context.Context) ([]*domain.Policy, error) {
	apiKey := config.E.YouthApiKey

	reqURL := fmt.Sprintf("https://www.youthcenter.go.kr/opi/empSprtList.do?openApiVcntId=%s&pageIndex=1&display=20", url.QueryEscape(apiKey))

	body, err := u.doGetRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	var items []domain.YouthPolicyItem

	var xmlResp domain.YouthPolicyXMLResponse
	if xmlErr := xml.Unmarshal(body, &xmlResp); xmlErr == nil && len(xmlResp.PolicyList) > 0 {
		items = xmlResp.PolicyList
	} else {
		var jsonResp domain.YouthPolicyJSONResponse
		if jsonErr := json.Unmarshal(body, &jsonResp); jsonErr == nil && len(jsonResp.EmpsInfo.Emp) > 0 {
			items = jsonResp.EmpsInfo.Emp
		}
	}

	var result []*domain.Policy
	now := time.Now()

	for _, item := range items {
		title := strings.TrimSpace(item.PolyBizSjnNm)
		if title == "" {
			continue
		}
		desc := strings.TrimSpace(item.PolyItcnCn)
		if desc == "" {
			desc = title
		}

		sDate, eDate := parseDateRange(item.RqstPrdCn)
		if sDate.IsZero() {
			sDate = now
		}
		if eDate.IsZero() {
			eDate = now.AddDate(1, 0, 0)
		}

		regionCode := parseRegionCode(item.PolyBizSecd, item.CnsgNtiPrdCn)

		result = append(result, &domain.Policy{
			Title:       title,
			Description: desc,
			RegionCode:  regionCode,
			StartDate:   sDate,
			EndDate:     eDate,
			Address:     strings.TrimSpace(item.CnsgNtiPrdCn),
			Latitude:    0.0,
			Longitude:   0.0,
		})
	}

	return result, nil
}

// FetchWelfarePolicies fetches welfare services from 한국사회보장정보원 (중앙부처복지서비스)
func (u *policyUseCase) FetchWelfarePolicies(ctx context.Context) ([]*domain.Policy, error) {
	items, err := u.fetchWelfareList(ctx)
	if err != nil {
		return nil, err
	}

	var (
		mu     sync.Mutex
		result = make([]*domain.Policy, 0, len(items))
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(1)

	for _, item := range items {
		time.Sleep(100 * time.Millisecond)
		item := item

		g.Go(func() error {
			exist, _ := u.policyRepository.FindByTitle(ctx, item.ServNm)
			if exist != nil {
				return nil
			}
			policy, err := u.fetchWelfareDetail(ctx, item.ServID)
			if err != nil {
				log.Printf("welfare detail %s: %v", item.ServID, err)
				return nil
			}

			mu.Lock()
			result = append(result, policy)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (u *policyUseCase) fetchWelfareDetail(
	ctx context.Context,
	servID string,
) (*domain.Policy, error) {

	params := url.Values{}
	params.Set("serviceKey", config.E.WelfareApiKey)
	params.Set("callTp", "D")
	params.Set("servId", servID)

	reqURL := fmt.Sprintf(
		"https://apis.data.go.kr/B554287/NationalWelfareInformationsV001/NationalWelfaredetailedV001?%s",
		params.Encode(),
	)

	body, err := u.doGetRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	var resp domain.WelfareDetailResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	item := resp.WelfareDetailItem

	var desc strings.Builder

	desc.WriteString(strings.TrimSpace(item.WlfareInfoOutlCn))

	if item.AlwServCn != "" {
		desc.WriteString("\n\n지원내용\n")
		desc.WriteString(item.AlwServCn)
	}

	if item.TgtrDtlCn != "" {
		desc.WriteString("\n\n지원대상\n")
		desc.WriteString(item.TgtrDtlCn)
	}

	if item.SlctCritCn != "" {
		desc.WriteString("\n\n선정기준\n")
		desc.WriteString(item.SlctCritCn)
	}

	if len(item.ApplmetList) > 0 {
		desc.WriteString("\n\n신청방법")
		for _, v := range item.ApplmetList {
			desc.WriteString("\n- ")
			desc.WriteString(v.ServSeDetailNm)
		}
	}

	address := item.JurMnofNm
	if item.RprsCtadr != "" {
		address = item.RprsCtadr
	}

	return &domain.Policy{
		Title:       strings.TrimSpace(item.ServNm),
		Description: desc.String(),
		RegionCode:  0, // 중앙부처
		StartDate:   time.Now(),
		EndDate:     time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
		Address:     address,
		Latitude:    0,
		Longitude:   0,
	}, nil
}

func (u *policyUseCase) fetchWelfareList(ctx context.Context) ([]domain.WelfareItem, error) {
	params := url.Values{}
	params.Set("serviceKey", config.E.WelfareApiKey)
	params.Set("callTp", "L")
	params.Set("pageNo", "1")
	params.Set("numOfRows", "100")
	params.Set("srchKeyCode", "001")

	reqURL := fmt.Sprintf(
		"https://apis.data.go.kr/B554287/NationalWelfareInformationsV001/NationalWelfarelistV001?%s",
		params.Encode(),
	)

	body, err := u.doGetRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("request welfare list: %w", err)
	}

	// XML 우선
	var xmlResp domain.WelfareXMLResponse
	if err := xml.Unmarshal(body, &xmlResp); err == nil {
		if xmlResp.ResultCode != "" && xmlResp.ResultCode != "0" {
			return nil, fmt.Errorf(
				"api error: %s (%s)",
				xmlResp.ResultMessage,
				xmlResp.ResultCode,
			)
		}

		if len(xmlResp.ServList) > 0 {
			return xmlResp.ServList, nil
		}
	}

	// JSON fallback
	var jsonResp domain.WelfareJSONResponse
	if err := json.Unmarshal(body, &jsonResp); err == nil {
		if jsonResp.WantedList.ResultCode != "" &&
			jsonResp.WantedList.ResultCode != "0" {
			return nil, fmt.Errorf(
				"api error: %s (%s)",
				jsonResp.WantedList.ResultMessage,
				jsonResp.WantedList.ResultCode,
			)
		}

		return jsonResp.WantedList.ServList, nil
	}

	return nil, fmt.Errorf("failed to decode welfare response")
}

// FetchVolunteerEvents fetches volunteer activities and events from 행정안전부 봉사참여정보서비스
func (u *policyUseCase) FetchVolunteerEvents(ctx context.Context) ([]*domain.Policy, error) {
	list, err := u.fetchVolunteerAreaList(ctx)
	if err != nil {
		return nil, err
	}

	var (
		mu     sync.Mutex
		result = make([]*domain.Policy, 0, len(list))
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(1)

	for _, item := range list {
		time.Sleep(100 * time.Millisecond)
		item := item

		g.Go(func() error {
			exist, _ := u.policyRepository.FindByTitle(ctx, item.ProgrmSj)
			if exist != nil {
				return nil
			}
			policy, err := u.fetchVolunteerDetail(ctx, item.ProgrmRegistNo)
			if err != nil {
				log.Printf("volunteer detail %s: %v", item.ProgrmRegistNo, err)
				return nil
			}

			mu.Lock()
			result = append(result, policy)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (u *policyUseCase) fetchVolunteerAreaList(ctx context.Context) ([]domain.VolunteerAreaItem, error) {
	params := url.Values{}
	params.Set("serviceKey", config.E.VolunteerApiKey)
	params.Set("pageNo", "1")
	params.Set("numOfRows", "50")
	params.Set("schSido", "6110000") // 서울

	reqURL := fmt.Sprintf(
		"https://apis.data.go.kr/1741000/volunteerPartcptnService/getVltrAreaList?%s",
		params.Encode(),
	)

	fmt.Println(reqURL)

	body, err := u.doGetRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("request volunteer area list: %w", err)
	}

	var resp domain.VolunteerAreaResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode volunteer area list: %w", err)
	}

	if resp.Header.ResultCode != "00" {
		return nil, fmt.Errorf("api error: %s (%s)",
			resp.Header.ResultMsg,
			resp.Header.ResultCode,
		)
	}

	return resp.Body.Items.Item, nil
}

func (u *policyUseCase) fetchVolunteerDetail(
	ctx context.Context,
	registNo string,
) (*domain.Policy, error) {

	params := url.Values{}
	params.Set("serviceKey", config.E.VolunteerApiKey)
	params.Set("pageNo", "1")
	params.Set("numOfRows", "1")
	params.Set("progrmRegistNo", registNo)

	reqURL := fmt.Sprintf(
		"https://apis.data.go.kr/1741000/volunteerPartcptnService/getVltrPartcptnItem?%s",
		params.Encode(),
	)

	body, err := u.doGetRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	var resp domain.VolunteerDetailResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	item := resp.Body.Items.Item

	start := parseYYYYMMDD(item.ProgrmBgnde)
	end := parseYYYYMMDD(item.ProgrmEndde)

	lat, _ := strconv.ParseFloat(item.AreaLalo1, 64)
	lng, _ := strconv.ParseFloat(item.AreaLalo2, 64)

	desc := fmt.Sprintf(
		`%s

주관기관 : %s
모집기관 : %s
모집상태 : %s
모집인원 : %d명
신청인원 : %d명
활동요일 : %s`,
		strings.TrimSpace(item.ProgrmCn),
		item.NanmmbyNm,
		item.MnnstNm,
		item.ProgrmSttusSe,
		item.RcritNmpr,
		item.AppTotal,
		item.ActWkdy,
	)

	return &domain.Policy{
		Title:       item.ProgrmSj,
		Description: desc,
		RegionCode:  parseRegionCode(item.SidoCd, item.GugunCd),
		StartDate:   start,
		EndDate:     end,
		Address: strings.TrimSpace(
			item.AreaAddress1 + " " +
				item.AreaAddress2 + " " +
				item.AreaAddress3,
		),
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

// FetchMaternityPolicies fetches maternity & childcare support status (도/시 출산장려/양육비 지원현황)
func (u *policyUseCase) FetchMaternityPolicies(ctx context.Context) ([]*domain.Policy, error) {
	apiKey := config.E.MaternityApiKey

	reqURL := fmt.Sprintf("http://apis.data.go.kr/1741000/MaternityChildcareService/getMaternityChildcareList?serviceKey=%s&ctpvNm=%s&pageNo=1&numOfRows=20", url.QueryEscape(apiKey), url.QueryEscape("서울특별시"))

	body, err := u.doGetRequest(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	var items []domain.MaternityItem

	var xmlResp domain.MaternityXMLResponse
	if xmlErr := xml.Unmarshal(body, &xmlResp); xmlErr == nil && len(xmlResp.Body.Items.ItemList) > 0 {
		items = xmlResp.Body.Items.ItemList
	} else {
		var jsonResp domain.MaternityJSONResponse
		if jsonErr := json.Unmarshal(body, &jsonResp); jsonErr == nil && len(jsonResp.Response.Body.Items.Item) > 0 {
			items = jsonResp.Response.Body.Items.Item
		}
	}

	var result []*domain.Policy
	now := time.Now()

	for _, item := range items {
		title := strings.TrimSpace(item.Bznm)
		if title == "" {
			title = strings.TrimSpace(item.ServNm)
		}
		if title == "" {
			continue
		}

		desc := strings.TrimSpace(item.DetlCn)
		if desc == "" {
			desc = strings.TrimSpace(item.ServDgst)
		}
		if desc == "" {
			desc = title
		}

		addr := strings.TrimSpace(fmt.Sprintf("%s %s", item.CtpvNm, item.SggNm))
		regionCode := parseRegionCode("", addr)

		sDate := parseYYYYMMDD(item.Bgnde)
		if sDate.IsZero() {
			sDate = now
		}
		eDate := parseYYYYMMDD(item.Endde)
		if eDate.IsZero() {
			eDate = now.AddDate(1, 0, 0)
		}

		result = append(result, &domain.Policy{
			Title:       title,
			Description: desc,
			RegionCode:  regionCode,
			StartDate:   sDate,
			EndDate:     eDate,
			Address:     addr,
			Latitude:    0.0,
			Longitude:   0.0,
		})
	}

	return result, nil
}

// Helper methods for HTTP, Regional Filtering and Data Parsing
func (u *policyUseCase) doGetRequest(ctx context.Context, reqURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/xml, application/json, text/xml, */*")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code: %d\nbody: %s",
			resp.StatusCode,
			string(body),
		)
	}

	return body, nil

	return io.ReadAll(resp.Body)
}

func isSeoulOrNational(p *domain.Policy) bool {
	// Central / National policy
	if p.RegionCode == 0 {
		if isNonSeoulText(p.Address) || isNonSeoulText(p.Title) {
			return false
		}
		return true
	}
	// Seoul Region Code (11000)
	if p.RegionCode == 11000 {
		return true
	}
	// Address containing Seoul
	if strings.Contains(p.Address, "서울") {
		return true
	}
	return false
}

func isNonSeoulText(text string) bool {
	nonSeoulKeywords := []string{
		"부산", "대구", "인천", "광주", "대전", "울산", "세종",
		"경기", "강원", "충북", "충남", "전북", "전남", "경북", "경남", "제주",
		"경상", "전라", "충청",
	}
	for _, kw := range nonSeoulKeywords {
		if strings.Contains(text, kw) && !strings.Contains(text, "서울") {
			return true
		}
	}
	return false
}

func parseYYYYMMDD(s string) time.Time {
	s = strings.TrimSpace(s)
	if len(s) < 8 {
		return time.Time{}
	}
	re := regexp.MustCompile(`\d{8}`)
	match := re.FindString(s)
	if match == "" {
		return time.Time{}
	}

	t, err := time.Parse("20060102", match)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseDateRange(rangeStr string) (time.Time, time.Time) {
	if rangeStr == "" {
		return time.Time{}, time.Time{}
	}

	re := regexp.MustCompile(`\d{4}[.-/]?\d{2}[.-/]?\d{2}`)
	matches := re.FindAllString(rangeStr, -1)

	var startDate, endDate time.Time

	if len(matches) >= 1 {
		startDate = cleanAndParseDate(matches[0])
	}
	if len(matches) >= 2 {
		endDate = cleanAndParseDate(matches[1])
	}

	return startDate, endDate
}

func cleanAndParseDate(dateStr string) time.Time {
	cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(dateStr, "")
	if len(cleaned) == 8 {
		t, err := time.Parse("20060102", cleaned)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}

func parseRegionCode(codeStr, addressStr string) int {
	if strings.Contains(addressStr, "서울") {
		return 11000
	}
	if codeStr != "" {
		cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(codeStr, "")
		if val, err := strconv.Atoi(cleaned); err == nil {
			if val == 11000 || val == 11 || val == 6110000 {
				return 11000
			}
		}
	}
	return 0
}
