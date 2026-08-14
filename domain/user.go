package domain

import (
	"context"
	"time"
)

type User struct {
	ID         ID
	Name       string
	Email      string
	Password   string
	Age        int
	RegionCode int
	Gender     Gender
	IsStudent  bool
	IsYouth    bool
	IsPregnant bool
	IsBusiness bool
	IsDisabled bool
	IsForeign  bool
	Interests  []Interest
	CreatedAt  time.Time
}

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

type Interest string

const (
	InterestEmployment    Interest = "EMPLOYMENT"
	InterestHousing       Interest = "HOUSING"
	InterestEducation     Interest = "EDUCATION"
	InterestWelfare       Interest = "WELFARE"
	InterestPregnancy     Interest = "PREGNANCY"
	InterestCulture       Interest = "CULTURE"
	InterestEnvironment   Interest = "ENVIRONMENT"
	InterestParticipation Interest = "PARTICIPATION"
)

type UserRepository interface {
	Create(c context.Context, user *User) (*User, error)
	Update(c context.Context, user *User) (*User, error)
	FindAll(c context.Context) ([]*User, error)
	FindByEmail(c context.Context, email string) (*User, error)
	FindByID(c context.Context, id *ID) (*User, error)
}
