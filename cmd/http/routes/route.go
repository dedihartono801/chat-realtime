package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetupRoutes(
	app fiber.Router,

) {
	app.Get("/", monitor.New())
}
