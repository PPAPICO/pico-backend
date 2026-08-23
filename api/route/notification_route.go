package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/pico-backend/api/handler"
	"github.com/janghanul090801/pico-backend/api/middleware"
	"github.com/janghanul090801/pico-backend/domain"
)

func NewNotificationRoute(app fiber.Router, service domain.NotificationUseCase) {
	protected := app.Group("/protected")
	protected.Use(middleware.JwtMiddleware)
	protected.Get("/", handler.GetNotifications(service))
	protected.Patch("/:id", handler.MarkNotificationAsRead(service))
	protected.Delete("/:id", handler.DeleteNotification(service))
	protected.Get("/unread-count", handler.GetUnreadNotificationsCount(service))
}