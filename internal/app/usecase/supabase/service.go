package supabase

import (
	"errors"

	"github.com/dedihartono801/chat-realtime/internal/app/repository"
	"github.com/dedihartono801/chat-realtime/internal/entity"
	"github.com/dedihartono801/chat-realtime/pkg/customstatus"
	"github.com/dedihartono801/chat-realtime/pkg/dto"
	"github.com/dedihartono801/chat-realtime/pkg/validator"
)

type Service interface {
	SaveMessages(req *dto.SaveMessageDto) (int, error)
}

type service struct {
	messageRepository repository.MessageRepository
	validator         validator.Validator
}

func NewSupabaseService(
	messageRepository repository.MessageRepository,
	validator validator.Validator,
) Service {
	return &service{
		messageRepository: messageRepository,
		validator:         validator,
	}
}

func (s *service) SaveMessages(req *dto.SaveMessageDto) (int, error) {
	if err := s.validator.Validate(req); err != nil {
		return customstatus.ErrBadRequest.Code, errors.New(customstatus.ErrBadRequest.Message)
	}

	message := &entity.Message{
		UserChatID:  req.UserChatID,
		MessageText: req.MessageText,
	}
	err := s.messageRepository.CreateMessage(message)
	if err != nil {
		return customstatus.ErrInternalServerError.Code, errors.New(err.Error())
	}
	return customstatus.StatusCreated.Code, nil

}
