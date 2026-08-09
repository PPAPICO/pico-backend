package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/janghanul090801/pico-backend/api/route"
	"github.com/janghanul090801/pico-backend/config"
	"github.com/janghanul090801/pico-backend/external/volunteer"
	"github.com/janghanul090801/pico-backend/external/welfare"
	"github.com/janghanul090801/pico-backend/external/youth"
	"github.com/janghanul090801/pico-backend/infra/database"
	"github.com/janghanul090801/pico-backend/infra/repository"
	"github.com/janghanul090801/pico-backend/internal/httpclient"
	"github.com/janghanul090801/pico-backend/usecase"

	_ "github.com/janghanul090801/pico-backend/docs"
)

// @title          PICO Backend API
// @version        1.0
// @description    PICO Backend API
// @host			localhost:8000
// @BasePath		/api
func main() {
	config.NewEnv()

	app := fiber.New(fiber.Config{
		AppName:      "Fiber Ent Clean Architecture",
		ServerHeader: "Fiber",
	})

	// Use global middlewares.
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(limiter.New(limiter.Config{
		Max: 100,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(&fiber.Map{
				"status":  "fail",
				"message": "You have requested too many in a single time-frame! Please wait another minute!",
			})
		},
	}))
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())

	api := app.Group("/api")

	app.Get("/swagger/*", swaggo.HandlerDefault)

	app.Get("/docs/*", swaggo.HandlerDefault)

	client, err := database.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	if err = client.Schema.Create(context.Background()); err != nil {
		panic(err)
	}

	timeout := time.Duration(config.E.ContextTimeout) * time.Second

	httpClient := httpclient.NewClient(&http.Client{
		Timeout: timeout,
	})

	// repository
	userRepository := repository.NewUserRepository(client)
	policyRepository := repository.NewPolicyRepository(client)
	policyMatchRepository := repository.NewPolicyMatchRepository(client)

	// external client
	youthClient := youth.NewClient(httpClient, config.E.YouthApiKey)
	welfareClient := welfare.NewClient(httpClient, policyRepository, config.E.WelfareApiKey)
	volunteerClient := volunteer.NewClient(httpClient, policyRepository, config.E.VolunteerApiKey)

	// usecase
	profileUseCase := usecase.NewProfileUseCase(userRepository, timeout)
	authUseCase := usecase.NewAuthUseCase(userRepository, timeout)
	policyUseCase := usecase.NewPolicyUseCase(policyRepository, policyMatchRepository, youthClient, welfareClient, volunteerClient, timeout)

	// router
	route.NewLoginRouter(api.Group("/login"), authUseCase)
	route.NewProfileRouter(api.Group("/profile"), profileUseCase)
	route.NewRefreshTokenRouter(api.Group("/refresh"), authUseCase)
	route.NewSignupRouter(api.Group("/signup"), authUseCase)
	route.NewPolicyRouter(api.Group("/policy"), policyUseCase, profileUseCase)

	app.All("*", func(c fiber.Ctx) error {
		notFoundErr := fmt.Errorf(
			"route '%s' does not exist in this API",
			c.OriginalURL(),
		)

		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"status": "error",
			"error":  notFoundErr.Error(),
		})
	})

	log.Fatal(app.Listen(config.E.ServerAddress))
}
