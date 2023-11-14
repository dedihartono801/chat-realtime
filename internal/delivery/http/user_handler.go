package http

import (
	"github.com/dedihartono801/chat-realtime/internal/app/usecase/user"
	"github.com/dedihartono801/chat-realtime/pkg/dto"
	"github.com/dedihartono801/chat-realtime/pkg/helpers"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Registration(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type userHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) UserHandler {
	return &userHandler{service}
}

func (h *userHandler) Registration(ctx *fiber.Ctx) error {
	userDto := new(dto.UserRegistrationDto)
	if err := ctx.BodyParser(userDto); err != nil {
		return helpers.CustomResponse(ctx, nil, err.Error(), 500)
	}
	user, statusCode, err := h.service.Register(userDto)

	if err != nil {
		return helpers.CustomResponse(ctx, nil, err.Error(), statusCode)
	}
	return helpers.CustomResponse(ctx, user, nil, statusCode)

}

func (h *userHandler) Login(ctx *fiber.Ctx) error {
	loginDto := new(dto.LoginDto)
	if err := ctx.BodyParser(loginDto); err != nil {
		return helpers.CustomResponse(ctx, nil, err.Error(), 500)
	}
	user, statusCode, err := h.service.Login(loginDto)

	if err != nil {
		return helpers.CustomResponse(ctx, nil, err.Error(), statusCode)
	}
	return helpers.CustomResponse(ctx, user, nil, statusCode)

}
