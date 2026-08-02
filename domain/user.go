package domain

import (
	"context"
	"time"
)

type User struct {
	ID          ID
	Name        string
	Email       string
	Password    string
	Age         int
	Region      string
	Gender      Gender
	IsStudent   bool
	IsYouth     bool
	IsPregnant  *bool
	IsBusiness  *bool
	Interests   []Interest
	IsDisabled  bool
	Nationality Nationality
	CreatedAt   time.Time
}

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

type Nationality string

const (
	NationalityDomestic Nationality = "DOMESTIC"
	NationalityForeign  Nationality = "FOREIGN"
)

type Interest string

const (
	InterestCulture        Interest = "CULTURE"
	InterestWelfare        Interest = "WELFARE"
	InterestEnvironment    Interest = "ENVIRONMENT"
	InterestEmployment     Interest = "EMPLOYMENT"
	InterestEducation      Interest = "EDUCATION"
	InterestHousing        Interest = "HOUSING"
	InterestStartup        Interest = "STARTUP"
	InterestEconomy        Interest = "ECONOMY"
	InterestHealthcare     Interest = "HEALTHCARE"
	InterestTransportation Interest = "TRANSPORTATION"
)

type UserRepository interface {
	Save(c context.Context, user *User) (*User, error)
	FindAll(c context.Context) ([]*User, error)
	FindByEmail(c context.Context, email string) (*User, error)
	FindByID(c context.Context, id *ID) (*User, error)
}
