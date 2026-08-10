package domain

import (
	"context"
)

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type SignupRequest struct {
	Name       string     `form:"name" binding:"required"`
	Email      string     `form:"email" binding:"required,email"`
	Password   string     `form:"password" binding:"required"`
	Age        int        `form:"age" binding:"required"`
	RegionCode int        `form:"region" binding:"required"`
	Gender     Gender     `form:"gender" binding:"required"`
	IsStudent  bool       `form:"is_student" binding:"required"`
	IsYouth    bool       `form:"is_youth" binding:"required"`
	IsPregnant bool       `form:"is_pregnant" binding:"required"`
	IsBusiness bool       `form:"is_business" binding:"required"`
	IsDisabled bool       `json:"is_disabled" binding:"required"`
	IsForeign  bool       `json:"is_foreign" binding:"required"`
	Interests  []Interest `form:"interests" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `form:"refreshToken" binding:"required"`
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
