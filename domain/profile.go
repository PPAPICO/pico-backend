package domain

import (
	"context"
	"time"
)

type Profile struct {
	ID         ID         `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Age        int        `json:"age"`
	RegionCode int        `json:"region"`
	Gender     Gender     `json:"gender"`
	IsStudent  bool       `json:"is_student"`
	IsYouth    bool       `json:"is_youth"`
	IsPregnant bool       `json:"is_pregnant"`
	IsBusiness bool       `json:"is_business"`
	IsDisabled bool       `json:"is_disabled"`
	IsForeign  bool       `json:"is_foreign"`
	Interests  []Interest `json:"interests"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ProfileUseCase interface {
	GetProfileByID(c context.Context, userID *ID) (*Profile, error)
}
