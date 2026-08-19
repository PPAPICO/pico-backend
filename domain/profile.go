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

type ProfileUpdateRequest struct {
	Name       string     `json:"name" form:"name" binding:"required"`
	Password   string     `json:"password" form:"password" binding:"required"`
	Age        int        `json:"age" form:"age" binding:"required"`
	RegionCode int        `json:"region" form:"region" binding:"required"`
	Gender     Gender     `json:"gender" form:"gender" binding:"required"`
	IsStudent  bool       `json:"is_student" form:"is_student" binding:"required"`
	IsYouth    bool       `json:"is_youth" form:"is_youth" binding:"required"`
	IsPregnant bool       `json:"is_pregnant" form:"is_pregnant" binding:"required"`
	IsBusiness bool       `json:"is_business" form:"is_business" binding:"required"`
	IsDisabled bool       `json:"is_disabled" form:"is_disabled" binding:"required"`
	IsForeign  bool       `json:"is_foreign" form:"is_foreign" binding:"required"`
	Interests  []Interest `json:"interests" form:"interests" binding:"required"`
}

type ProfileUseCase interface {
	GetByID(c context.Context, userID *ID) (*Profile, error)
	Update(c context.Context, ID *ID, params *ProfileUpdateRequest) (*Profile, error)
}
