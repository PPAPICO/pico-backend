package youth

import "encoding/xml"

// 온통청년 청년정책 API Structs

type Item struct {
	PlcyNo      string `json:"plcyNo" xml:"plcyNo"`
	PlcyNm      string `json:"plcyNm" xml:"plcyNm"`
	PlcyKywdNm  string `json:"plcyKywdNm" xml:"plcyKywdNm"`
	PlcyExplnCn string `json:"plcyExplnCn" xml:"plcyExplnCn"`
	LclsfNm     string `json:"lclsfNm" xml:"lclsfNm"`
	MclsfNm     string `json:"mclsfNm" xml:"mclsfNm"`
	PlcySprtCn  string `json:"plcySprtCn" xml:"plcySprtCn"`

	SprvsnInstCd   string `json:"sprvsnInstCd" xml:"sprvsnInstCd"`
	SprvsnInstCdNm string `json:"sprvsnInstCdNm" xml:"sprvsnInstCdNm"`
	OperInstCdNm   string `json:"operInstCdNm" xml:"operInstCdNm"`

	BizPrdBgngYmd string `json:"bizPrdBgngYmd" xml:"bizPrdBgngYmd"`
	BizPrdEndYmd  string `json:"bizPrdEndYmd" xml:"bizPrdEndYmd"`
	BizPrdEtcCn   string `json:"bizPrdEtcCn" xml:"bizPrdEtcCn"`

	AplyYmd     string `json:"aplyYmd" xml:"aplyYmd"`
	AplyUrlAddr string `json:"aplyUrlAddr" xml:"aplyUrlAddr"`

	ZipCd string `json:"zipCd" xml:"zipCd"`

	SprtTrgtMinAge string `json:"sprtTrgtMinAge" xml:"sprtTrgtMinAge"`
	SprtTrgtMaxAge string `json:"sprtTrgtMaxAge" xml:"sprtTrgtMaxAge"`

	RefUrlAddr1 string `json:"refUrlAddr1" xml:"refUrlAddr1"`
	RefUrlAddr2 string `json:"refUrlAddr2" xml:"refUrlAddr2"`
}

type response struct {
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	Result        struct {
		Pagging struct {
			TotCount int `json:"totCount"`
			PageNum  int `json:"pageNum"`
			PageSize int `json:"pageSize"`
		} `json:"pagging"`

		YouthPolicyList []Item `json:"youthPolicyList"`
	} `json:"result"`
}

type XMLResponse struct {
	XMLName    xml.Name `xml:"youthPolicyList"`
	PolicyList []Item   `xml:"youthPolicy"`
}
