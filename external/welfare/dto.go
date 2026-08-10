package welfare

import "encoding/xml"

// 한국사회보장정보원 복지서비스정보 (중앙부처복지서비스 API Structs)

type XMLResponse struct {
	XMLName       xml.Name `xml:"wantedList"`
	ResultCode    string   `xml:"resultCode"`
	ResultMessage string   `xml:"resultMessage"`
	TotalCount    int      `xml:"totalCount"`
	PageNo        int      `xml:"pageNo"`
	NumOfRows     int      `xml:"numOfRows"`
	ServList      []Item   `xml:"servList"`
}

type JSONResponse struct {
	WantedList struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
		TotalCount    int    `json:"totalCount,string"`
		PageNo        int    `json:"pageNo,string"`
		NumOfRows     int    `json:"numOfRows,string"`
		ServList      []Item `json:"servList"`
	} `json:"wantedList"`
}

type Body struct {
	Items []Item `xml:"items"`
}

type Item struct {
	InqNum            string `xml:"inqNum" json:"inqNum"`
	ServID            string `xml:"servId" json:"servId"`
	ServNm            string `xml:"servNm" json:"servNm"`
	ServDgst          string `xml:"servDgst" json:"servDgst"`
	JurMnofNm         string `xml:"jurMnofNm" json:"jurMnofNm"`
	JurOrgNm          string `xml:"jurOrgNm" json:"jurOrgNm"`
	ServDtlLink       string `xml:"servDtlLink" json:"servDtlLink"`
	RprsCtadr         string `xml:"rprsCtadr" json:"rprsCtadr"`
	SprtCycNm         string `xml:"sprtCycNm" json:"sprtCycNm"`
	SrvPvsnNm         string `xml:"srvPvsnNm" json:"srvPvsnNm"`
	LifeArray         string `xml:"lifeArray" json:"lifeArray"`
	TrgterIndvdlArray string `xml:"trgterIndvdlArray" json:"trgterIndvdlArray"`
	IntrsThemaArray   string `xml:"intrsThemaArray" json:"intrsThemaArray"`
	OnapPsbltYn       string `xml:"onapPsbltYn" json:"onapPsbltYn"`
}

type DetailResponse struct {
	XMLName xml.Name `xml:"wantedDtl"`
	DetailItem
}

type DetailItem struct {
	ServID           string `xml:"servId"`
	ServNm           string `xml:"servNm"`
	JurMnofNm        string `xml:"jurMnofNm"`
	RprsCtadr        string `xml:"rprsCtadr"`
	TgtrDtlCn        string `xml:"tgtrDtlCn"`
	SlctCritCn       string `xml:"slctCritCn"`
	AlwServCn        string `xml:"alwServCn"`
	WlfareInfoOutlCn string `xml:"wlfareInfoOutlCn"`

	ApplmetList []LinkItem `xml:"applmetList"`
}

type LinkItem struct {
	ServSeCode       string `xml:"servSeCode"`
	ServSeDetailNm   string `xml:"servSeDetailNm"`
	ServSeDetailLink string `xml:"servSeDetailLink"`
}
