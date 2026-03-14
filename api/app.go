package api

import (
	"EtuEDT-Go/domain"
	"errors"
	"log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func StartWebApp() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) {
				return c.Status(e.Code).JSON(domain.ErrorResponse{Error: e.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "internal server error"})
		},
	})
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(recover.New())

	prometheus := fiberprometheus.New("EtuEDT-Back")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/monitor", monitor.New())
	app.Get("/openapi", func(c *fiber.Ctx) error {
		return c.SendFile("./openapi.yaml")
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"path": "v3",
		})
	})
	v2Router(app.Group("/v2"))
	v3Router(app.Group("/v3"))

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
