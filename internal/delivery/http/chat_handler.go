package http

import (
	"fmt"

	"github.com/dedihartono801/chat-realtime/internal/app/usecase/chat"
	"github.com/dedihartono801/chat-realtime/pkg/dto"
	"github.com/dedihartono801/chat-realtime/pkg/helpers"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler interface {
	SendMessage(ctx *fiber.Ctx) error
	FetchMessage(ctx *fiber.Ctx) error
	SearchMessage(ctx *fiber.Ctx) error
}

type chatHandler struct {
	service chat.Service
}

func NewChatHandler(service chat.Service) ChatHandler {
	return &chatHandler{service}
}

func (h *chatHandler) SendMessage(ctx *fiber.Ctx) error {
	messageDto := new(dto.SendMessageDto)
	if err := ctx.BodyParser(messageDto); err != nil {
		return helpers.CustomResponse(ctx, nil, err.Error(), 500)
	}

	// Assuming ctx.Locals("ID") returns an interface{}
	val := ctx.Locals("ID")

	statusCode, err := h.service.SendMessage(messageDto, val.(int64), ctx.Params("to"))

	if err != nil {
		return helpers.CustomResponse(ctx, nil, err.Error(), statusCode)
	}
	return helpers.CustomResponse(ctx, nil, "success", statusCode)

}

func (h *chatHandler) FetchMessage(ctx *fiber.Ctx) error {
	// Assuming ctx.Locals("ID") returns an interface{}
	val := ctx.Locals("ID")

	data, statusCode, err := h.service.FetchMessage(val.(int64), ctx.Params("to"))

	if err != nil {
		return helpers.CustomResponse(ctx, nil, err.Error(), statusCode)
	}
	return helpers.CustomResponse(ctx, data, "success", statusCode)

}

func (h *chatHandler) SearchMessage(ctx *fiber.Ctx) error {
	// Assuming ctx.Locals("ID") returns an interface{}
	val := ctx.Locals("ID")
	message := ctx.Query("message")
	fmt.Println(message)

	data, statusCode, err := h.service.SearchMessage(val.(int64), message)

	if err != nil {
		return helpers.CustomResponse(ctx, nil, err.Error(), statusCode)
	}
	return helpers.CustomResponse(ctx, data, "success", statusCode)

}
