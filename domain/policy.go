package domain

import (
	"context"
	"encoding/xml"
	"strings"
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

type PolicyResponse struct {
	ID          ID        `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	RegionCode  int       `json:"region_code"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Address     string    `json:"address"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Status      Match     `json:"status"`
}

func (p *Policy) ToResponse(match *PolicyMatch) *PolicyResponse {
	status := MatchUNCERTAIN
	if match != nil {
		status = match.Status
	}
	return &PolicyResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		RegionCode:  p.RegionCode,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		Address:     p.Address,
		Latitude:    p.Latitude,
		Longitude:   p.Longitude,
		Status:      status,
	}
}

type Match int

const (
	MatchPOSSIBLE   Match = iota
	MatchUNCERTAIN  Match = iota
	MatchIMPOSSIBLE Match = iota
)

type XMLHeader struct {
	ResultCode string `xml:"resultCode" json:"resultCode"`
	ResultMsg  string `xml:"resultMsg" json:"resultMsg"`
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

func (p *Policy) IsSeoulOrNational() bool {
	if p.RegionCode == 11000 {
		return true
	}

	if p.RegionCode == 0 {
		return !isNonSeoulText(p.Address) &&
			!isNonSeoulText(p.Title)
	}

	return strings.Contains(p.Address, "서울")
}

func isNonSeoulText(text string) bool {
	nonSeoulKeywords := []string{
		"부산", "대구", "인천", "광주", "대전", "울산", "세종",
		"경기", "강원", "충북", "충남", "전북", "전남", "경북", "경남", "제주",
		"경상", "전라", "충청",
	}
	for _, kw := range nonSeoulKeywords {
		if strings.Contains(text, kw) && !strings.Contains(text, "서울") {
			return true
		}
	}
	return false
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
	ListByRegionCodeAndActive(c context.Context, regionCode int) ([]*Policy, error)
	GetFromApi(c context.Context) ([]*Policy, error)
	ListMatchesByUserID(c context.Context, userID *ID) ([]*PolicyMatch, error)
}
