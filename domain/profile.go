package domain

import (
	"context"
	"time"
)

type Profile struct {
	ID          ID          `json:"id"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Age         int         `json:"age"`
	RegionCode  int         `json:"region"`
	Gender      Gender      `json:"gender"`
	IsStudent   bool        `json:"is_student"`
	IsYouth     bool        `json:"is_youth"`
	IsPregnant  *bool       `json:"is_pregnant"`
	IsBusiness  *bool       `json:"is_business"`
	Interests   []Interest  `json:"interests"`
	IsDisabled  bool        `json:"is_disabled"`
	Nationality Nationality `json:"nationality"`
	CreatedAt   time.Time   `json:"created_at"`
}

type ProfileUseCase interface {
	GetProfileByID(c context.Context, userID *ID) (*Profile, error)
}
