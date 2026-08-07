package domain

import (
	"context"
	"encoding/xml"
	"time"
)

type Policy struct {
	ID          ID
	Title       string
	Description string
	RegionCode  int
	StartDate   time.Time
	EndDate     time.Time
	Address     string
	Latitude    float64
	Longitude   float64
}

type PolicyMatch struct {
	ID       ID
	PolicyID ID
	UserID   ID
	Status   Match
}

type Match int

const (
	POSSIBLE   Match = iota
	UNCERTAIN  Match = iota
	IMPOSSIBLE Match = iota
)

// 1. 온통청년 청년정책 API Structs
type YouthPolicyXMLResponse struct {
	XMLName    xml.Name          `xml:"empsInfo"`
	TotalCnt   string            `xml:"totalCnt"`
	PageIndex  string            `xml:"pageIndex"`
	PolicyList []YouthPolicyItem `xml:"emp"`
}

type YouthPolicyJSONResponse struct {
	EmpsInfo struct {
		Emp []YouthPolicyItem `json:"emp"`
	} `json:"empsInfo"`
}

type YouthPolicyItem struct {
	BizID        string `xml:"bizId" json:"bizId"`
	PolyBizSjnNm string `xml:"polyBizSjnNm" json:"polyBizSjnNm"`
	PolyItcnCn   string `xml:"polyItcnCn" json:"polyItcnCn"`
	PolyBizSecd  string `xml:"polyBizSecd" json:"polyBizSecd"`
	RqstPrdCn    string `xml:"rqstPrdCn" json:"rqstPrdCn"`
	CnsgNtiPrdCn string `xml:"cnsgNtiPrdCn" json:"cnsgNtiPrdCn"`
	BizPrdCn     string `xml:"bizPrdCn" json:"bizPrdCn"`
}

// 2. 한국사회보장정보원 복지서비스정보 (중앙부처복지서비스 API Structs)
type WelfareXMLResponse struct {
	XMLName       xml.Name      `xml:"wantedList"`
	ResultCode    string        `xml:"resultCode"`
	ResultMessage string        `xml:"resultMessage"`
	TotalCount    int           `xml:"totalCount"`
	PageNo        int           `xml:"pageNo"`
	NumOfRows     int           `xml:"numOfRows"`
	ServList      []WelfareItem `xml:"servList"`
}

type WelfareJSONResponse struct {
	WantedList struct {
		ResultCode    string        `json:"resultCode"`
		ResultMessage string        `json:"resultMessage"`
		TotalCount    int           `json:"totalCount,string"`
		PageNo        int           `json:"pageNo,string"`
		NumOfRows     int           `json:"numOfRows,string"`
		ServList      []WelfareItem `json:"servList"`
	} `json:"wantedList"`
}

type XMLHeader struct {
	ResultCode string `xml:"resultCode" json:"resultCode"`
	ResultMsg  string `xml:"resultMsg" json:"resultMsg"`
}

type WelfareBody struct {
	Items WelfareItems `xml:"items"`
}

type WelfareItems struct {
	ItemList []WelfareItem `xml:"item"`
}

type WelfareItem struct {
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

type WelfareDetailResponse struct {
	XMLName xml.Name `xml:"wantedDtl"`
	WelfareDetailItem
}

type WelfareDetailItem struct {
	ServID           string `xml:"servId"`
	ServNm           string `xml:"servNm"`
	JurMnofNm        string `xml:"jurMnofNm"`
	RprsCtadr        string `xml:"rprsCtadr"`
	TgtrDtlCn        string `xml:"tgtrDtlCn"`
	SlctCritCn       string `xml:"slctCritCn"`
	AlwServCn        string `xml:"alwServCn"`
	WlfareInfoOutlCn string `xml:"wlfareInfoOutlCn"`

	ApplmetList []WelfareLinkItem `xml:"applmetList"`
}

type WelfareLinkItem struct {
	ServSeCode       string `xml:"servSeCode"`
	ServSeDetailNm   string `xml:"servSeDetailNm"`
	ServSeDetailLink string `xml:"servSeDetailLink"`
}

// 3. 행정안전부 봉사참여정보서비스 API Structs
type VolunteerAreaResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode" json:"resultCode"`
		ResultMsg  string `xml:"resultMsg" json:"resultMsg"`
	} `xml:"header" json:"header"`

	Body struct {
		Items struct {
			Item []VolunteerAreaItem `xml:"item" json:"item"`
		} `xml:"items" json:"items"`

		PageNo     int `xml:"pageNo" json:"pageNo"`
		NumOfRows  int `xml:"numOfRows" json:"numOfRows"`
		TotalCount int `xml:"totalCount" json:"totalCount"`
	} `xml:"body" json:"body"`
}

type VolunteerAreaItem struct {
	GugunCd        string `xml:"gugunCd" json:"gugunCd"`
	SidoCd         string `xml:"sidoCd" json:"sidoCd"`
	ProgrmRegistNo string `xml:"progrmRegistNo" json:"progrmRegistNo"`
	ProgrmSj       string `xml:"progrmSj" json:"progrmSj"`
	ProgrmBgnde    string `xml:"progrmBgnde" json:"progrmBgnde"`
	ProgrmEndde    string `xml:"progrmEndde" json:"progrmEndde"`
	ProgrmSttusSe  string `xml:"progrmSttusSe" json:"progrmSttusSe"`
	NanmmbyNm      string `xml:"nanmmbyNm" json:"nanmmbyNm"`
}

type VolunteerDetailResponse struct {
	XMLName xml.Name `xml:"response"`

	Header VolunteerHeader `xml:"header"`
	Body   struct {
		Items struct {
			Item VolunteerDetailItem `xml:"item"`
		} `xml:"items"`

		PageNo     int `xml:"pageNo"`
		NumOfRows  int `xml:"numOfRows"`
		TotalCount int `xml:"totalCount"`
	} `xml:"body"`
}

type VolunteerHeader struct {
	ResultCode string `xml:"resultCode"`
	ResultMsg  string `xml:"resultMsg"`
}

type VolunteerDetailItem struct {
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

// 4. 도/시 출산장려/양육비 지원현황 API Structs
type MaternityXMLResponse struct {
	XMLName xml.Name      `xml:"response"`
	Header  XMLHeader     `xml:"header"`
	Body    MaternityBody `xml:"body"`
}

type MaternityJSONResponse struct {
	Response struct {
		Header XMLHeader `json:"header"`
		Body   struct {
			Items struct {
				Item []MaternityItem `json:"item"`
			} `json:"items"`
		} `json:"body"`
	} `json:"response"`
}

type MaternityBody struct {
	Items MaternityItems `xml:"items"`
}

type MaternityItems struct {
	ItemList []MaternityItem `xml:"item"`
}

type MaternityItem struct {
	Bznm     string `xml:"bznm" json:"bznm"`
	ServNm   string `xml:"servNm" json:"servNm"`
	DetlCn   string `xml:"detlCn" json:"detlCn"`
	ServDgst string `xml:"servDgst" json:"servDgst"`
	CtpvNm   string `xml:"ctpvNm" json:"ctpvNm"`
	SggNm    string `xml:"sggNm" json:"sggNm"`
	Bgnde    string `xml:"bgnde" json:"bgnde"`
	Endde    string `xml:"endde" json:"endde"`
}

type PolicyRepository interface {
	FindAll(c context.Context) ([]*Policy, error)
	FindByID(c context.Context, id *ID) (*Policy, error)
	FindByTitle(c context.Context, title string) (*Policy, error)
	FindAllByRegionCode(c context.Context, regionCode int) ([]*Policy, error)
	FindAllByRegionCodeAndActive(c context.Context, regionCode int) ([]*Policy, error)
	Create(c context.Context, policy *Policy) (*Policy, error)
	Delete(c context.Context, id *ID) error
}

type PolicyMatchRepository interface {
	FindAllByUserID(c context.Context, userID *ID) ([]*PolicyMatch, error)
	FindAllByUserIDAndStatus(c context.Context, userID *ID, status Match) ([]*PolicyMatch, error)
	Create(c context.Context, policyMatch *PolicyMatch) (*PolicyMatch, error)
	Update(c context.Context, id *ID, status Match) (*PolicyMatch, error)
}

type PolicyUseCase interface {
	GetByID(c context.Context, id *ID) (*Policy, error)
	ListRegionCodeAndActive(c context.Context, regionCode int) ([]*Policy, error)
	GetFromApi(c context.Context) ([]*Policy, error)
	GetMatchesByUserID(c context.Context, userID *ID) ([]*Policy, error)
}
