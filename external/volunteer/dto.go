package volunteer

import "encoding/xml"

// 행정안전부 봉사참여정보서비스 API Structs

type AreaResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode" json:"resultCode"`
		ResultMsg  string `xml:"resultMsg" json:"resultMsg"`
	} `xml:"header" json:"header"`

	Body struct {
		Items struct {
			Item []AreaItem `xml:"item" json:"item"`
		} `xml:"items" json:"items"`

		PageNo     int `xml:"pageNo" json:"pageNo"`
		NumOfRows  int `xml:"numOfRows" json:"numOfRows"`
		TotalCount int `xml:"totalCount" json:"totalCount"`
	} `xml:"body" json:"body"`
}

type AreaItem struct {
	GugunCd        string `xml:"gugunCd" json:"gugunCd"`
	SidoCd         string `xml:"sidoCd" json:"sidoCd"`
	ProgrmRegistNo string `xml:"progrmRegistNo" json:"progrmRegistNo"`
	ProgrmSj       string `xml:"progrmSj" json:"progrmSj"`
	ProgrmBgnde    string `xml:"progrmBgnde" json:"progrmBgnde"`
	ProgrmEndde    string `xml:"progrmEndde" json:"progrmEndde"`
	ProgrmSttusSe  string `xml:"progrmSttusSe" json:"progrmSttusSe"`
	NanmmbyNm      string `xml:"nanmmbyNm" json:"nanmmbyNm"`
}

type DetailResponse struct {
	XMLName xml.Name `xml:"response"`

	Header Header `xml:"header"`
	Body   struct {
		Items struct {
			Item DetailItem `xml:"item"`
		} `xml:"items"`

		PageNo     int `xml:"pageNo"`
		NumOfRows  int `xml:"numOfRows"`
		TotalCount int `xml:"totalCount"`
	} `xml:"body"`
}

type Header struct {
	ResultCode string `xml:"resultCode"`
	ResultMsg  string `xml:"resultMsg"`
}

type DetailItem struct {
	ProgrmRegistNo string `xml:"progrmRegistNo"`
	ProgrmSj       string `xml:"progrmSj"`
	ProgrmCn       string `xml:"progrmCn"`

	ProgrmSttusSe string `xml:"progrmSttusSe"`

	ProgrmBgnde string `xml:"progrmBgnde"`
	ProgrmEndde string `xml:"progrmEndde"`

	ActBeginTm int `xml:"actBeginTm"`
	ActEndTm   int `xml:"actEndTm"`

	NoticeBgnde string `xml:"noticeBgnde"`
	NoticeEndde string `xml:"noticeEndde"`

	RcritNmpr int `xml:"rcritNmpr"`
	AppTotal  int `xml:"appTotal"`

	ActWkdy string `xml:"actWkdy"`

	SrvcClCode string `xml:"srvcClCode"`

	AdultPosblAt  string `xml:"adultPosblAt"`
	YngbgsPosblAt string `xml:"yngbgsPosblAt"`
	GrpPosblAt    string `xml:"grpPosblAt"`
	PbsvntPosblAt string `xml:"pbsvntPosblAt"`
	FamilyPosblAt string `xml:"familyPosblAt"`

	MnnstNm   string `xml:"mnnstNm"`
	NanmmbyNm string `xml:"nanmmbyNm"`

	ActPlace      string `xml:"actPlace"`
	NanmmbyNmAdmn string `xml:"nanmmbyNmAdmn"`
	Telno         string `xml:"telno"`
	Fxnum         string `xml:"fxnum"`
	PostAdres     string `xml:"postAdres"`
	Email         string `xml:"email"`

	SidoCd  string `xml:"sidoCd"`
	GugunCd string `xml:"gugunCd"`

	AreaAddress1 string `xml:"areaAddress1"`
	AreaAddress2 string `xml:"areaAddress2"`
	AreaAddress3 string `xml:"areaAddress3"`

	AreaLalo1 string `xml:"areaLalo1"` // latitude
	AreaLalo2 string `xml:"areaLalo2"` // longitude
	AreaLalo3 string `xml:"areaLalo3"`
}
