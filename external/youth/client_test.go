package youth

//
//import (
//	"context"
//	"encoding/xml"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestFetchYouthPolicies(t *testing.T) {
//	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		xmlResp := `<?xml version="1.0" encoding="UTF-8"?>
//<empsInfo>
//	<totalCnt>1</totalCnt>
//	<emp>
//		<bizId>Y0001</bizId>
//		<polyBizSjnNm>청년 월세 지원</polyBizSjnNm>
//		<polyItcnCn>월세 지원 사업</polyItcnCn>
//		<polyBizSecd>11000</polyBizSecd>
//		<rqstPrdCn>2026.01.01 ~ 2026.12.31</rqstPrdCn>
//		<cnsgNtiPrdCn>서울특별시 강남구</cnsgNtiPrdCn>
//	</emp>
//</empsInfo>`
//		w.Header().Set("Content-Type", "application/xml")
//		w.WriteHeader(http.StatusOK)
//		_, _ = w.Write([]byte(xmlResp))
//	}))
//	defer mockServer.Close()
//
//	uc := &policyUseCase{
//		httpClient: mockServer.Client(),
//	}
//
//	body, err := uc.doGetRequest(context.Background(), mockServer.URL)
//	assert.NoError(t, err)
//
//	var xmlResp domain.YouthPolicyXMLResponse
//	err = xml.Unmarshal(body, &xmlResp)
//	assert.NoError(t, err)
//	assert.Equal(t, 1, len(xmlResp.PolicyList))
//	assert.Equal(t, "청년 월세 지원", xmlResp.PolicyList[0].PolyBizSjnNm)
//}
