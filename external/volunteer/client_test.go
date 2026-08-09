package volunteer

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain/mocks"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTransportHTTPClient(transport roundTripFunc) *httpclient.Client {
	return httpclient.NewClient(
		&http.Client{
			Transport: transport,
		},
	)
}

func newPolicyRepository() domain.PolicyRepository {
	policyRepository := new(mocks.PolicyRepository{})
	policyRepository.On(
		"FindByTitle",
		mock.Anything,
		mock.AnythingOfType("string"),
	).Return(nil, nil)

	return policyRepository
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

		client := &Client{
			httpClient:       newTransportHTTPClient(mockTransport),
			policyRepository: newPolicyRepository(),
		}

		policies, err := client.Fetch(context.Background())
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

		client := &Client{
			httpClient:       newTransportHTTPClient(mockTransport),
			policyRepository: newPolicyRepository(),
		}

		policies, err := client.Fetch(context.Background())
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

		client := &Client{
			httpClient:       newTransportHTTPClient(mockTransport),
			policyRepository: newPolicyRepository(),
		}

		policies, err := client.Fetch(context.Background())
		assert.Error(t, err)
		assert.Nil(t, policies)
		assert.Contains(t, err.Error(), "unexpected status code: 502")
	})
}
