package usecase

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPolicyRepository struct {
	mock.Mock
}

func (m *mockPolicyRepository) FindAll(c context.Context) ([]*domain.Policy, error) {
	args := m.Called(c)
	return args.Get(0).([]*domain.Policy), args.Error(1)
}

func (m *mockPolicyRepository) FindByID(c context.Context, id *domain.ID) (*domain.Policy, error) {
	args := m.Called(c, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Policy), args.Error(1)
}

func (m *mockPolicyRepository) FindAllByRegionCode(c context.Context, regionCode int) ([]*domain.Policy, error) {
	args := m.Called(c, regionCode)
	return args.Get(0).([]*domain.Policy), args.Error(1)
}

func (m *mockPolicyRepository) FindAllByRegionCodeAndActive(c context.Context, regionCode int) ([]*domain.Policy, error) {
	args := m.Called(c, regionCode)
	return args.Get(0).([]*domain.Policy), args.Error(1)
}

func (m *mockPolicyRepository) Create(c context.Context, policy *domain.Policy) (*domain.Policy, error) {
	args := m.Called(c, policy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Policy), args.Error(1)
}

func (m *mockPolicyRepository) Delete(c context.Context, id *domain.ID) error {
	args := m.Called(c, id)
	return args.Error(0)
}

type mockPolicyMatchRepository struct {
	mock.Mock
}

func (m *mockPolicyMatchRepository) FindAllByUserID(c context.Context, userID *domain.ID) ([]*domain.PolicyMatch, error) {
	args := m.Called(c, userID)
	return args.Get(0).([]*domain.PolicyMatch), args.Error(1)
}

func (m *mockPolicyMatchRepository) FindAllByUserIDAndStatus(c context.Context, userID *domain.ID, status domain.Match) ([]*domain.PolicyMatch, error) {
	args := m.Called(c, userID, status)
	return args.Get(0).([]*domain.PolicyMatch), args.Error(1)
}

func (m *mockPolicyMatchRepository) Create(c context.Context, policyMatch *domain.PolicyMatch) (*domain.PolicyMatch, error) {
	args := m.Called(c, policyMatch)
	return args.Get(0).(*domain.PolicyMatch), args.Error(1)
}

func (m *mockPolicyMatchRepository) Update(c context.Context, id *domain.ID, status domain.Match) (*domain.PolicyMatch, error) {
	args := m.Called(c, id, status)
	return args.Get(0).(*domain.PolicyMatch), args.Error(1)
}

func TestParseDateHelpers(t *testing.T) {
	t.Run("parseYYYYMMDD", func(t *testing.T) {
		dateStr := "20260801"
		parsed := parseYYYYMMDD(dateStr)
		assert.Equal(t, 2026, parsed.Year())
		assert.Equal(t, time.Month(8), parsed.Month())
		assert.Equal(t, 1, parsed.Day())
	})

	t.Run("parseDateRange", func(t *testing.T) {
		rangeStr := "2026.08.01 ~ 2026.12.31"
		sDate, eDate := parseDateRange(rangeStr)
		assert.Equal(t, 2026, sDate.Year())
		assert.Equal(t, 2026, eDate.Year())
		assert.Equal(t, time.Month(12), eDate.Month())
	})
}

func TestParseRegionCode(t *testing.T) {
	assert.Equal(t, 11000, parseRegionCode("", "서울특별시 강남구"))
	assert.Equal(t, 11000, parseRegionCode("11000", ""))
	assert.Equal(t, 0, parseRegionCode("", "경기도 성남시"))
	assert.Equal(t, 0, parseRegionCode("", "부산광역시 해운대구"))
}

func TestIsSeoulOrNationalFilter(t *testing.T) {
	seoulPolicy := &domain.Policy{Title: "서울 청년 수당", RegionCode: 11000, Address: "서울특별시"}
	nationalPolicy := &domain.Policy{Title: "중앙 부처 복지", RegionCode: 0, Address: "보건복지부"}
	gyeonggiPolicy := &domain.Policy{Title: "경기 청년 지원", RegionCode: 41000, Address: "경기도 성남시"}
	busanPolicy := &domain.Policy{Title: "부산 출산 지원", RegionCode: 26000, Address: "부산광역시"}

	assert.True(t, isSeoulOrNational(seoulPolicy))
	assert.True(t, isSeoulOrNational(nationalPolicy))
	assert.False(t, isSeoulOrNational(gyeonggiPolicy))
	assert.False(t, isSeoulOrNational(busanPolicy))
}

func TestFetchYouthPolicies(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xmlResp := `<?xml version="1.0" encoding="UTF-8"?>
<empsInfo>
	<totalCnt>1</totalCnt>
	<emp>
		<bizId>Y0001</bizId>
		<polyBizSjnNm>청년 월세 지원</polyBizSjnNm>
		<polyItcnCn>월세 지원 사업</polyItcnCn>
		<polyBizSecd>11000</polyBizSecd>
		<rqstPrdCn>2026.01.01 ~ 2026.12.31</rqstPrdCn>
		<cnsgNtiPrdCn>서울특별시 강남구</cnsgNtiPrdCn>
	</emp>
</empsInfo>`
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(xmlResp))
	}))
	defer mockServer.Close()

	uc := &policyUseCase{
		httpClient: mockServer.Client(),
	}

	body, err := uc.doGetRequest(context.Background(), mockServer.URL)
	assert.NoError(t, err)

	var xmlResp domain.YouthPolicyXMLResponse
	err = xml.Unmarshal(body, &xmlResp)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(xmlResp.PolicyList))
	assert.Equal(t, "청년 월세 지원", xmlResp.PolicyList[0].PolyBizSjnNm)
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFetchWelfarePolicies(t *testing.T) {
	t.Run("Success - XML response", func(t *testing.T) {
		mockTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			urlStr := req.URL.String()

			if strings.Contains(urlStr, "NationalWelfarelistV001") {
				xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<wantedList>
	<resultCode>0</resultCode>
	<resultMessage>SUCCESS</resultMessage>
	<totalCount>1</totalCount>
	<servList>
		<servId>WLF00000001</servId>
		<servNm>청년 월세 복지 지원</servNm>
	</servList>
</wantedList>`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(xmlData)),
					Header:     make(http.Header),
				}, nil
			}

			if strings.Contains(urlStr, "NationalWelfaredetailedV001") {
				xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<wantedDtl>
	<servId>WLF00000001</servId>
	<servNm>청년 월세 복지 지원</servNm>
	<wlfareInfoOutlCn>월세 지원 개요</wlfareInfoOutlCn>
	<alwServCn>월 20만원 지원</alwServCn>
	<tgtrDtlCn>무주택 청년</tgtrDtlCn>
	<slctCritCn>중위소득 60% 이하</slctCritCn>
	<applmetList>
		<servSeDetailNm>복지로 온라인 신청</servSeDetailNm>
	</applmetList>
	<jurMnofNm>보건복지부</jurMnofNm>
</wantedDtl>`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(xmlData)),
					Header:     make(http.Header),
				}, nil
			}

			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("Not Found")),
			}, nil
		})

		uc := &policyUseCase{
			httpClient: &http.Client{Transport: mockTransport},
		}

		policies, err := uc.FetchWelfarePolicies(context.Background())
		assert.NoError(t, err)
		assert.Len(t, policies, 1)
		assert.Equal(t, "청년 월세 복지 지원", policies[0].Title)
		assert.Contains(t, policies[0].Description, "월세 지원 개요")
		assert.Contains(t, policies[0].Description, "월 20만원 지원")
		assert.Contains(t, policies[0].Description, "무주택 청년")
		assert.Contains(t, policies[0].Description, "복지로 온라인 신청")
		assert.Equal(t, "보건복지부", policies[0].Address)
		assert.Equal(t, 0, policies[0].RegionCode)
	})

	t.Run("Success - JSON response fallback", func(t *testing.T) {
		mockTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			urlStr := req.URL.String()

			if strings.Contains(urlStr, "NationalWelfarelistV001") {
				jsonData := `{
	"wantedList": {
		"resultCode": "0",
		"resultMessage": "SUCCESS",
		"servList": [
			{
				"servId": "WLF00000002",
				"servNm": "청년 자산 형성 지원"
			}
		]
	}
}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(jsonData)),
					Header:     make(http.Header),
				}, nil
			}

			if strings.Contains(urlStr, "NationalWelfaredetailedV001") {
				xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<wantedDtl>
	<servId>WLF00000002</servId>
	<servNm>청년 자산 형성 지원</servNm>
	<wlfareInfoOutlCn>자산 형성 지원 사업</wlfareInfoOutlCn>
	<rprsCtadr>서울특별시 종로구 세종대로 209</rprsCtadr>
</wantedDtl>`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(xmlData)),
					Header:     make(http.Header),
				}, nil
			}

			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("Not Found")),
			}, nil
		})

		uc := &policyUseCase{
			httpClient: &http.Client{Transport: mockTransport},
		}

		policies, err := uc.FetchWelfarePolicies(context.Background())
		assert.NoError(t, err)
		assert.Len(t, policies, 1)
		assert.Equal(t, "청년 자산 형성 지원", policies[0].Title)
		assert.Equal(t, "서울특별시 종로구 세종대로 209", policies[0].Address)
	})

	t.Run("Failure - API Error response", func(t *testing.T) {
		mockTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<wantedList>
	<resultCode>99</resultCode>
	<resultMessage>AUTHENTICATION ERROR</resultMessage>
</wantedList>`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(xmlData)),
				Header:     make(http.Header),
			}, nil
		})

		uc := &policyUseCase{
			httpClient: &http.Client{Transport: mockTransport},
		}

		policies, err := uc.FetchWelfarePolicies(context.Background())
		assert.Error(t, err)
		assert.Nil(t, policies)
		assert.Contains(t, err.Error(), "AUTHENTICATION ERROR")
	})

	t.Run("Failure - HTTP status error", func(t *testing.T) {
		mockTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("Internal Server Error")),
			}, nil
		})

		uc := &policyUseCase{
			httpClient: &http.Client{Transport: mockTransport},
		}

		policies, err := uc.FetchWelfarePolicies(context.Background())
		assert.Error(t, err)
		assert.Nil(t, policies)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})
}

func TestFetchVolunteerEvents(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			urlStr := req.URL.String()

			if strings.Contains(urlStr, "getVltrAreaList") {
				xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<response>
	<header>
		<resultCode>00</resultCode>
		<resultMsg>NORMAL SERVICE.</resultMsg>
	</header>
	<body>
		<items>
			<item>
				<progrmRegistNo>V123456</progrmRegistNo>
				<progrmSj>서울 환경 정화 봉사활동</progrmSj>
				<sidoCd>6110000</sidoCd>
				<gugunCd>11000</gugunCd>
			</item>
		</items>
		<pageNo>1</pageNo>
		<numOfRows>50</numOfRows>
		<totalCount>1</totalCount>
	</body>
</response>`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(xmlData)),
					Header:     make(http.Header),
				}, nil
			}

			if strings.Contains(urlStr, "getVltrPartcptnItem") {
				xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<response>
	<header>
		<resultCode>00</resultCode>
		<resultMsg>NORMAL SERVICE.</resultMsg>
	</header>
	<body>
		<items>
			<item>
				<progrmRegistNo>V123456</progrmRegistNo>
				<progrmSj>서울 환경 정화 봉사활동</progrmSj>
				<progrmCn>한강 플로깅 봉사 활동입니다.</progrmCn>
				<progrmSttusSe>모집중</progrmSttusSe>
				<progrmBgnde>20260301</progrmBgnde>
				<progrmEndde>20260331</progrmEndde>
				<rcritNmpr>20</rcritNmpr>
				<appTotal>5</appTotal>
				<actWkdy>주말</actWkdy>
				<mnnstNm>서울시 봉사센터</mnnstNm>
				<nanmmbyNm>서울 환경재단</nanmmbyNm>
				<sidoCd>6110000</sidoCd>
				<gugunCd>11000</gugunCd>
				<areaAddress1>서울특별시</areaAddress1>
				<areaAddress2>영등포구</areaAddress2>
				<areaAddress3>여의도동</areaAddress3>
				<areaLalo1>37.5283</areaLalo1>
				<areaLalo2>126.9294</areaLalo2>
			</item>
		</items>
	</body>
</response>`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(xmlData)),
					Header:     make(http.Header),
				}, nil
			}

			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("Not Found")),
			}, nil
		})

		uc := &policyUseCase{
			httpClient: &http.Client{Transport: mockTransport},
		}

		policies, err := uc.FetchVolunteerEvents(context.Background())
		assert.NoError(t, err)
		assert.Len(t, policies, 1)
		assert.Equal(t, "서울 환경 정화 봉사활동", policies[0].Title)
		assert.Contains(t, policies[0].Description, "한강 플로깅 봉사 활동입니다.")
		assert.Contains(t, policies[0].Description, "모집인원 : 20명")
		assert.Equal(t, 11000, policies[0].RegionCode)
		assert.Equal(t, 2026, policies[0].StartDate.Year())
		assert.Equal(t, time.Month(3), policies[0].StartDate.Month())
		assert.Equal(t, 1, policies[0].StartDate.Day())
		assert.Equal(t, 2026, policies[0].EndDate.Year())
		assert.Equal(t, time.Month(3), policies[0].EndDate.Month())
		assert.Equal(t, 31, policies[0].EndDate.Day())
		assert.Equal(t, "서울특별시 영등포구 여의도동", policies[0].Address)
		assert.InDelta(t, 37.5283, policies[0].Latitude, 0.0001)
		assert.InDelta(t, 126.9294, policies[0].Longitude, 0.0001)
	})

	t.Run("Failure - API Error response", func(t *testing.T) {
		mockTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<response>
	<header>
		<resultCode>99</resultCode>
		<resultMsg>SERVICE ERROR</resultMsg>
	</header>
</response>`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(xmlData)),
				Header:     make(http.Header),
			}, nil
		})

		uc := &policyUseCase{
			httpClient: &http.Client{Transport: mockTransport},
		}

		policies, err := uc.FetchVolunteerEvents(context.Background())
		assert.Error(t, err)
		assert.Nil(t, policies)
		assert.Contains(t, err.Error(), "SERVICE ERROR (99)")
	})

	t.Run("Failure - HTTP Status Error", func(t *testing.T) {
		mockTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Body:       io.NopCloser(strings.NewReader("Bad Gateway")),
			}, nil
		})

		uc := &policyUseCase{
			httpClient: &http.Client{Transport: mockTransport},
		}

		policies, err := uc.FetchVolunteerEvents(context.Background())
		assert.Error(t, err)
		assert.Nil(t, policies)
		assert.Contains(t, err.Error(), "unexpected status code: 502")
	})
}

