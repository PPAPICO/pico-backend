package usecase

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
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
