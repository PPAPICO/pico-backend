package welfare

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

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

		client := &Client{
			httpClient:       newTransportHTTPClient(mockTransport),
			policyRepository: newPolicyRepository(),
		}

		policies, err := client.Fetch(context.Background())
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

		client := &Client{
			httpClient:       newTransportHTTPClient(mockTransport),
			policyRepository: newPolicyRepository(),
		}

		policies, err := client.Fetch(context.Background())
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

		client := &Client{
			httpClient:       newTransportHTTPClient(mockTransport),
			policyRepository: newPolicyRepository(),
		}

		policies, err := client.Fetch(context.Background())
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

		client := &Client{
			httpClient:       newTransportHTTPClient(mockTransport),
			policyRepository: newPolicyRepository(),
		}

		policies, err := client.Fetch(context.Background())
		assert.Error(t, err)
		assert.Nil(t, policies)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})
}
