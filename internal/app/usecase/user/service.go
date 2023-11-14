package user

import (
	"errors"
	"time"

	"github.com/dedihartono801/chat-realtime/cmd/http/middleware"
	"github.com/dedihartono801/chat-realtime/internal/app/repository"
	"github.com/dedihartono801/chat-realtime/internal/entity"
	"github.com/dedihartono801/chat-realtime/pkg/customstatus"
	"github.com/dedihartono801/chat-realtime/pkg/dto"
	"github.com/dedihartono801/chat-realtime/pkg/helpers"
	"github.com/dedihartono801/chat-realtime/pkg/validator"
)

type Service interface {
	Register(req *dto.UserRegistrationDto) (*entity.User, int, error)
	Login(req *dto.LoginDto) (*dto.LoginResponse, int, error)
}

type service struct {
	userRepository repository.UserRepository
	validator      validator.Validator
}

func NewUserService(
	userRepository repository.UserRepository,
	validator validator.Validator,
) Service {
	return &service{
		userRepository: userRepository,
		validator:      validator,
	}
}

func (s *service) Register(req *dto.UserRegistrationDto) (*entity.User, int, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, customstatus.ErrBadRequest.Code, errors.New(customstatus.ErrBadRequest.Message)
	}

	_, err := s.userRepository.GetUserByUsername(req.Username)
	if err == nil {
		return nil, customstatus.ErrConflictUsername.Code, errors.New(customstatus.ErrConflictUsername.Message)
	}

	if req.Password != "" {
		req.Password = helpers.EncryptPassword(req.Password)
	}

	user := &entity.User{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
	}

	dt, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, customstatus.ErrInternalServerError.Code, errors.New(err.Error())
	}

	return dt, customstatus.StatusCreated.Code, nil
}

func (s *service) Login(req *dto.LoginDto) (*dto.LoginResponse, int, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, customstatus.ErrBadRequest.Code, errors.New(customstatus.ErrBadRequest.Message)
	}

	user, err := s.userRepository.GetUserByUsername(req.Username)
	if err != nil {
		return nil, customstatus.ErrEmailNotFound.Code, errors.New(customstatus.ErrEmailNotFound.Message)
	}

	if user.Password != helpers.EncryptPassword(req.Password) {
		return nil, customstatus.ErrPasswordWrong.Code, errors.New(customstatus.ErrPasswordWrong.Message)
	}
	expirationTime := time.Now().Add(time.Hour * time.Duration(24))
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return nil, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}
	responseParams := dto.LoginResponse{
		Token:     token,
		ExpiredAt: expirationTime.Format(time.RFC3339),
	}
	return &responseParams, customstatus.StatusOk.Code, nil
}
