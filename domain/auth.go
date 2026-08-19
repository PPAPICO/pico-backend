package domain

import (
	"context"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type SignupRequest struct {
	Name       string     `json:"name" form:"name" binding:"required"`
	Email      string     `json:"email" form:"email" binding:"required,email"`
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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" form:"refreshToken" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthUseCase interface {
	Register(c context.Context, name, email, password string, age, regionCode int, gender Gender, isStudent, isYouth, isPregnant, isBusiness, isDisabled, isForeign bool, interest []Interest) (*User, error)
	Login(c context.Context, email, password string) (*User, error)
	CreateAccessAndRefreshToken(c context.Context, user *User) (string, string, error)
	ExtractUserFromRefreshToken(c context.Context, requestToken string) (*User, error)
}
