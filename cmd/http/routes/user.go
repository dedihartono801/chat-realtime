package routes

import (
	"github.com/dedihartono801/chat-realtime/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

// UserRouter is the Router for GoFiber App
func UserRouter(app fiber.Router, userHandler http.UserHandler) {
	usersRoute := app.Group("/users")
	usersRoute.Post("/registration", userHandler.Registration)
	usersRoute.Post("/login", userHandler.Login)
}
