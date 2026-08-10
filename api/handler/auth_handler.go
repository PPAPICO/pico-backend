package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/pico-backend/domain"
)

// Login
// @Summary      로그인
// @Description  로그인
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body  domain.LoginRequest  true  "로그인 요청"
// @Success      200  {object}  domain.AuthResponse
// @Failure      400  {object}  domain.ErrorResponse  "잘못된 요청"
// @Failure      401  {object}  domain.ErrorResponse  "인증 실패"
// @Failure      500  {object}  domain.ErrorResponse  "서버 오류"
// @Router       /login [post]
func Login(service domain.AuthUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()

		var request domain.LoginRequest

		err := c.Bind().Body(&request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		user, err := service.Login(ctx, request.Email, request.Password)
		if err != nil {
			if errInfo, ok := errors.AsType[domain.Error](err); ok {
				return c.Status(errInfo.StatusCode).JSON(domain.ErrorResponse{Message: err.Error()})
			}
		}

		accessToken, refreshToken, err := service.CreateAccessAndRefreshToken(ctx, user)
		if err != nil {
			if errInfo, ok := errors.AsType[domain.Error](err); ok {
				return c.Status(errInfo.StatusCode).JSON(domain.ErrorResponse{Message: err.Error()})
			}
		}

		response := domain.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return c.Status(http.StatusOK).JSON(response)
	}
}

// RefreshToken
// @Summary      토큰 갱신
// @Description  Refresh Token 으로 새로운 Access Token과 Refresh Token을 발급
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body  domain.RefreshTokenRequest  true  "토큰 갱신 요청"
// @Success      200  {object}  domain.AuthResponse
// @Failure      400  {object}  domain.ErrorResponse  "잘못된 요청"
// @Failure      401  {object}  domain.ErrorResponse  "유효하지 않은 Refresh Token"
// @Failure      500  {object}  domain.ErrorResponse  "서버 오류"
// @Router       /refresh [post]
func RefreshToken(service domain.AuthUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()

		var request domain.RefreshTokenRequest

		err := c.Bind().Body(&request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		user, err := service.ExtractUserFromRefreshToken(ctx, request.RefreshToken)
		if err != nil {
			if errInfo, ok := errors.AsType[domain.Error](err); ok {
				return c.Status(errInfo.StatusCode).JSON(domain.ErrorResponse{Message: err.Error()})
			}
		}

		accessToken, refreshToken, err := service.CreateAccessAndRefreshToken(ctx, user)
		if err != nil {
			if errInfo, ok := errors.AsType[domain.Error](err); ok {
				return c.Status(errInfo.StatusCode).JSON(domain.ErrorResponse{Message: err.Error()})
			}
		}

		response := domain.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return c.Status(http.StatusOK).JSON(response)
	}
}

// Signup
// @Summary      회원가입
// @Description  이름, 이메일, 비밀번호를 사용하여 회원가입
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body  domain.SignupRequest  true  "회원가입 요청"
// @Success      200  {object}  domain.AuthResponse
// @Failure      400  {object}  domain.ErrorResponse  "잘못된 요청"
// @Failure      409  {object}  domain.ErrorResponse  "이미 존재하는 사용자"
// @Failure      500  {object}  domain.ErrorResponse  "서버 오류"
// @Router       /signup [post]
func Signup(service domain.AuthUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()

		var request domain.SignupRequest

		err := c.Bind().Body(&request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		user, err := service.Register(
			ctx,
			request.Name,
			request.Email,
			request.Password,
			request.Age,
			request.RegionCode,
			request.Gender,
			request.IsStudent,
			request.IsYouth,
			request.IsPregnant,
			request.IsBusiness,
			request.IsDisabled,
			request.IsForeign,
			request.Interests,
		)
		if err != nil {
			return RespondError(c, err)
		}

		accessToken, refreshToken, err := service.CreateAccessAndRefreshToken(ctx, user)
		if err != nil {
			return RespondError(c, err)
		}

		response := domain.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return c.Status(http.StatusOK).JSON(response)
	}
}
