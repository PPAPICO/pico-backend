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
// @Router       /policy/protected/region [get]
func GetPoliciesInMyRegion(policyService domain.PolicyUseCase, profileService domain.ProfileUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		userID := c.Locals("id").(*domain.ID)

		profile, err := profileService.GetByID(ctx, userID)
		if err != nil {
			return RespondError(c, err)
		}

		policies, err := policyService.ListByRegionCodeAndActive(ctx, profile.RegionCode)
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

// GetAllPolicies
// @Summary      정책 조회
// @Description  정책 리스트 반환
// @Tags         Policy
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  domain.PolicyResponse
// @Failure      401  {object}  domain.ErrorResponse  "인증되지 않은 사용자"
// @Failure      500  {object}  domain.ErrorResponse  "서버 오류"
// @Router       /policy [get]
func GetAllPolicies(policyService domain.PolicyUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()

		policies, err := policyService.List(ctx)
		if err != nil {
			return RespondError(c, err)
		}

		return c.Status(http.StatusOK).JSON(policies)
	}
}

// AddFavoritePolicy
// @Summary      정책 즐겨찾기 등록
// @Description  특정 정책을 즐겨찾기에 등록합니다.
// @Tags         Policy
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "정책 ID"
// @Success      200  {object}  domain.Favorite
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /policy/protected/favorite/{id} [post]
func AddFavoritePolicy(favoriteUseCase domain.FavoriteUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		userID := c.Locals("id").(*domain.ID)
		policyID, err := domain.StringToID(c.Params("id"))
		if err != nil {
			return RespondError(c, domain.NewBadRequestError(err))
		}

		fav, err := favoriteUseCase.AddFavorite(ctx, userID, &policyID)
		if err != nil {
			return RespondError(c, err)
		}

		return c.Status(http.StatusOK).JSON(fav)
	}
}

// RemoveFavoritePolicy
// @Summary      정책 즐겨찾기 해제
// @Description  특정 정책을 즐겨찾기에서 제거합니다.
// @Tags         Policy
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "정책 ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /policy/protected/favorite/{id} [delete]
func RemoveFavoritePolicy(favoriteUseCase domain.FavoriteUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		userID := c.Locals("id").(*domain.ID)
		policyID, err := domain.StringToID(c.Params("id"))
		if err != nil {
			return RespondError(c, domain.NewBadRequestError(err))
		}

		err = favoriteUseCase.RemoveFavorite(ctx, userID, &policyID)
		if err != nil {
			return RespondError(c, err)
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
	}
}

// GetFavoritePolicies
// @Summary      즐겨찾기 정책 목록 조회
// @Description  로그인한 사용자의 즐겨찾기한 정책 목록을 조회합니다.
// @Tags         Policy
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   domain.Policy
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /policy/protected/favorites [get]
func GetFavoritePolicies(favoriteUseCase domain.FavoriteUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		userID := c.Locals("id").(*domain.ID)

		policies, err := favoriteUseCase.ListFavoritesByUserID(ctx, userID)
		if err != nil {
			return RespondError(c, err)
		}

		return c.Status(http.StatusOK).JSON(policies)
	}
}
