package routes

import (
	"github.com/dedihartono801/chat-realtime/cmd/http/middleware"
	"github.com/dedihartono801/chat-realtime/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

// UserRouter is the Router for GoFiber App
func ChatRouter(app fiber.Router, chatHandler http.ChatHandler) {
	chatRoute := app.Group("/chat", middleware.AuthUser)
	chatRoute.Put("/send-message/:to", chatHandler.SendMessage)
	chatRoute.Get("/search", chatHandler.SearchMessage)
	chatRoute.Get("/:to", chatHandler.FetchMessage)
}
