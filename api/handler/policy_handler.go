package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/internal/collections"
)

// GetPoliciesInMyRegion
// @Summary      내 지역 정책 조회
// @Description  현재 로그인한 사용자의 지역에 해당하는 정책 리스트 반환
// @Tags         Policy
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  domain.PolicyResponse
// @Failure      401  {object}  domain.ErrorResponse  "인증되지 않은 사용자"
// @Failure      500  {object}  domain.ErrorResponse  "서버 오류"
// @Router       /policy/protected [get]
func GetPoliciesInMyRegion(policyService domain.PolicyUseCase, profileService domain.ProfileUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		userID := c.Locals("id").(*domain.ID)

		user, err := profileService.GetProfileByID(ctx, userID)
		if err != nil {
			return RespondError(c, err)
		}

		policies, err := policyService.ListByRegionCodeAndActive(ctx, user.RegionCode)
		if err != nil {
			return RespondError(c, err)
		}

		matches, err := policyService.ListMatchesByUserID(ctx, userID)
		if err != nil {
			return RespondError(c, err)
		}

		policyResponses := collections.Map(policies, func(p *domain.Policy) *domain.PolicyResponse {
			return p.ToResponse(*collections.Find(matches, func(m *domain.PolicyMatch) bool {
				return m.PolicyID == p.ID
			}))
		})

		return c.Status(http.StatusOK).JSON(policyResponses)
	}
}
