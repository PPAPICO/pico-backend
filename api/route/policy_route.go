package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/pico-backend/api/handler"
	"github.com/janghanul090801/pico-backend/api/middleware"
	"github.com/janghanul090801/pico-backend/domain"
)

func NewPolicyRouter(app fiber.Router, policyService domain.PolicyUseCase, profileService domain.ProfileUseCase, favoriteService domain.FavoriteUseCase) {
	app.Get("/", handler.GetAllPolicies(policyService))

	// protected
	protected := app.Group("/protected")
	protected.Use(middleware.JwtMiddleware)
	protected.Get("/region", handler.GetPoliciesInMyRegion(policyService, profileService))
	protected.Get("/favorites", handler.GetFavoritePolicies(favoriteService, policyService))
	protected.Post("/favorite/:id", handler.AddFavoritePolicy(favoriteService))
	protected.Delete("/favorite/:id", handler.RemoveFavoritePolicy(favoriteService))
	protected.Post("/reclassify", handler.ReclassifyPolicies(policyService))
}
