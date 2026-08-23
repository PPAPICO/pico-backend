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
	Condition   PolicyCondition
}

type PolicyMatch struct {
	ID          ID
	PolicyID    ID
	UserID      ID
	Status      Match
	Probability *int
	Comment     *string
}

type PolicyResponse struct {
	ID          ID         `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	RegionCode  int        `json:"region_code"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Address     string     `json:"address"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	Status      Match      `json:"status"`
	Interests   []Interest `json:"interests"`
	SourceURL   string     `json:"source_url"`
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
		Interests:   p.Condition.Interests,
		SourceURL:   p.Condition.SourceURL,
	}
}

type Match int

const (
	MatchPOSSIBLE Match = iota
	MatchUNCERTAIN
	MatchIMPOSSIBLE
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
			!isNonSeoulText(p.Title) && !isNonSeoulText(p.Description)
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
	UpdateCondition(c context.Context, id *ID, condition PolicyCondition) error
}

type PolicyMatchRepository interface {
	FindAllByUserID(c context.Context, userID *ID) ([]*PolicyMatch, error)
	FindAllByUserIDAndStatus(c context.Context, userID *ID, status Match) ([]*PolicyMatch, error)
	FindByUserIDAndPolicyID(c context.Context, userID *ID, policyID *ID) (*PolicyMatch, error)
	Create(c context.Context, policyMatch *PolicyMatch) (*PolicyMatch, error)
	Update(c context.Context, id *ID, status Match, probability *int, comment *string) (*PolicyMatch, error)
}

type PolicyUseCase interface {
	GetByID(c context.Context, id *ID) (*Policy, error)
	ListByRegionCodeAndActive(c context.Context, regionCode int) ([]*Policy, error)
	GetFromApi(c context.Context) ([]*Policy, error)
	ListMatchesByUserID(c context.Context, userID *ID) ([]*PolicyMatch, error)
	List(c context.Context) ([]*Policy, error)
	SavePolicyMatches(c context.Context, user *User, policies []*Policy) ([]*PolicyMatch, error)
	UpdatePolicyMatchesForUser(c context.Context, userID *ID) ([]*PolicyMatch, error)
	ReclassifyAllPolicies(c context.Context) (int, error)
}

const MatchPrompt = "너는 사용자의 개인정보와 정책/봉사활동 정보를 비교하여, 해당 사용자가 해당 정책 또는 봉사활동의 대상자에 해당할 가능성을 판단하는 AI다.\n\n다음 두 가지 정보를 바탕으로 판단하라.\n\n[사용자 정보]\n%s\n\n[정책 또는 봉사활동 정보]\n%s\n\n판단 기준:\n\n1. 정책 또는 봉사활동의 지원 대상, 참여 조건, 자격 요건을 우선적으로 확인한다.\n2. 사용자의 나이, 성별, 지역, 국적, 소득, 직업, 학생 여부, 관심 분야 등 제공된 정보와 조건을 비교한다.\n3. 명시적인 자격 조건이 사용자 정보와 일치하면 높은 점수를 부여한다.\n4. 명시적인 자격 조건과 사용자 정보가 충돌하면 낮은 점수를 부여한다.\n5. 사용자 정보에 자격 판단에 필요한 정보가 없으면 해당 조건을 충족한다고 추측하지 말고 불확실성을 반영한다.\n6. 정책/봉사활동에 별도의 자격 제한이 없거나 대부분의 사람이 참여할 수 있는 경우에는 높은 점수를 부여할 수 있다.\n7. 제공된 정보만 사용하여 판단하며, 존재하지 않는 사용자 정보를 추측하거나 만들어내지 않는다.\n8. 확률은 0~99 사이의 정수 하나로 표현한다.\n9. 확률은 \"해당 사용자가 이 정책/봉사활동의 대상자 또는 참여 가능 대상일 가능성\"을 의미한다.\n10. 판단 근거는 핵심적인 이유를 한 문장으로 간결하게 작성한다.\n11. 반드시 아래 형식으로만 응답한다.\n\n출력 형식:\n{확률},{판단 근거}\n\n예시:\n85,사용자가 서울에 거주하고 만 19세 이상이라는 조건을 충족하므로 대상자일 가능성이 높습니다.\n42,지역 조건은 충족하지만 소득 조건을 판단할 사용자 정보가 없어 대상 여부가 불확실합니다.\n10,해당 정책은 만 65세 이상을 대상으로 하지만 사용자는 해당 연령 조건을 충족하지 않습니다.\n\n주의:\n\n* 확률은 반드시 0~99 사이의 정수여야 한다.\n* 확률 뒤에는 반드시 쉼표 하나만 사용한다.\n* 판단 근거에는 불필요한 설명이나 여러 문장을 넣지 않는다.\n* JSON, Markdown, 코드 블록, 접두사/접미사 등 다른 형식은 사용하지 않는다.\n* 최종 응답은 반드시 \"{정수},{한 문장}\" 형태여야 한다.\n"
const InterestPrompt = "너는 정책의 제목과 설명을 분석하여, 해당 정책이 어떤 분야에 속하는지 분류하는 AI다.\n\n[정책 정보]\n%s\n\n다음 분야 중 정책의 주요 목적과 내용을 기준으로 해당하는 분야를 하나 이상 선택하라.\n\n* EMPLOYMENT: 취업, 일자리, 고용, 창업, 직업훈련, 근로 지원 등\n* HOUSING: 주거, 주택, 전월세, 임대, 주거비, 주거환경 지원 등\n* EDUCATION: 교육, 학습, 장학금, 학교, 교육비, 역량교육 등\n* WELFARE: 복지, 생활비, 사회보장, 취약계층 지원, 생계 지원 등\n* PREGNANCY: 임신, 출산, 출산지원, 임산부, 육아, 산모·신생아 지원 등\n* CULTURE: 문화, 예술, 공연, 체육, 여가, 문화활동 지원 등\n* ENVIRONMENT: 환경보호, 기후, 에너지, 자원순환, 친환경 활동 등\n* PARTICIPATION: 봉사활동, 시민참여, 사회참여, 지역사회 활동, 공익활동 등\n\n분류 기준:\n\n1. 정책의 단순한 부수적 내용보다 정책의 핵심 목적과 지원 내용을 기준으로 판단한다.\n2. 여러 분야와 직접적으로 관련된 정책이라면 여러 분야를 선택할 수 있다.\n3. 단순히 관련 키워드가 등장한다는 이유만으로 해당 분야를 선택하지 않는다.\n4. 정책의 핵심 목적이 명확하게 특정 분야에 해당하면 해당 분야만 선택한다.\n5. 정책 설명만으로 특정 분야에 해당한다고 판단하기 어려우면 억지로 분류하지 않는다.\n6. 반드시 위 목록에 존재하는 분야만 사용한다.\n7. 분야 이름은 반드시 대문자 영문 코드로 출력한다.\n8. 선택한 분야는 쉼표(,)로만 구분한다.\n9. 별도의 설명, 이유, 문장, JSON, Markdown, 코드 블록은 출력하지 않는다.\n10. 가장 관련성이 높은 분야부터 순서대로 출력한다.\n\n출력 형식:\n{분야1},{분야2},{분야3}\n\n예시:\nEMPLOYMENT\nEMPLOYMENT,EDUCATION\nHOUSING,WELFARE\nPREGNANCY,WELFARE\nCULTURE,PARTICIPATION\nENVIRONMENT,PARTICIPATION\n\n반드시 위 출력 형식만 사용하라.\n"
