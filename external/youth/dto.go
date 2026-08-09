package youth

import "encoding/xml"

// 온통청년 청년정책 API Structs

type XMLResponse struct {
	XMLName    xml.Name `xml:"empsInfo"`
	TotalCnt   string   `xml:"totalCnt"`
	PageIndex  string   `xml:"pageIndex"`
	PolicyList []Item   `xml:"emp"`
}

type JSONResponse struct {
	EmpsInfo struct {
		Emp []Item `json:"emp"`
	} `json:"empsInfo"`
}

type Item struct {
	BizID        string `xml:"bizId" json:"bizId"`
	PolyBizSjnNm string `xml:"polyBizSjnNm" json:"polyBizSjnNm"`
	PolyItcnCn   string `xml:"polyItcnCn" json:"polyItcnCn"`
	PolyBizSecd  string `xml:"polyBizSecd" json:"polyBizSecd"`
	RqstPrdCn    string `xml:"rqstPrdCn" json:"rqstPrdCn"`
	CnsgNtiPrdCn string `xml:"cnsgNtiPrdCn" json:"cnsgNtiPrdCn"`
	BizPrdCn     string `xml:"bizPrdCn" json:"bizPrdCn"`
}
