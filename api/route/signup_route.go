package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/pico-backend/api/handler"
	"github.com/janghanul090801/pico-backend/domain"
)

func NewSignupRouter(app fiber.Router, service domain.AuthUseCase) {
	app.Post("/", handler.Signup(service))
}
