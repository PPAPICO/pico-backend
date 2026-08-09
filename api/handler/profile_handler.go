package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/pico-backend/domain"
)

// FetchProfile
// @Summary      내 프로필 조회
// @Description  현재 로그인한 사용자의 프로필 정보를 조회합니다.
// @Tags         Profile
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  domain.Profile
// @Failure      401  {object}  domain.ErrorResponse  "인증되지 않은 사용자"
// @Failure      404  {object}  domain.ErrorResponse  "프로필을 찾을 수 없음"
// @Failure      500  {object}  domain.ErrorResponse  "서버 오류"
// @Router       /profile/protected [get]
func FetchProfile(service domain.ProfileUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		userID := c.Locals("id").(*domain.ID)

		profile, err := service.GetProfileByID(ctx, userID)
		if err != nil {
			return RespondError(c, err)
		}

		return c.Status(http.StatusOK).JSON(profile)
	}
}
