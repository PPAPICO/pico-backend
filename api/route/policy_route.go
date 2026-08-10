package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/pico-backend/api/handler"
	"github.com/janghanul090801/pico-backend/api/middleware"
	"github.com/janghanul090801/pico-backend/domain"
)

func NewPolicyRouter(app fiber.Router, policyService domain.PolicyUseCase, profileService domain.ProfileUseCase) {
	// protected
	protected := app.Group("/protected")
	protected.Use(middleware.JwtMiddleware)
	protected.Get("/", handler.GetPoliciesInMyRegion(policyService, profileService))
}
