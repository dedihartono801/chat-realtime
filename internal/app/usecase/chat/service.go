package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/dedihartono801/chat-realtime/internal/app/queue/kafka"
	"github.com/dedihartono801/chat-realtime/internal/app/repository"
	"github.com/dedihartono801/chat-realtime/internal/entity"
	"github.com/dedihartono801/chat-realtime/pkg/config"
	"github.com/dedihartono801/chat-realtime/pkg/customstatus"
	"github.com/dedihartono801/chat-realtime/pkg/dto"
	"github.com/dedihartono801/chat-realtime/pkg/helpers"
	"github.com/dedihartono801/chat-realtime/pkg/validator"
)

type Service interface {
	SendMessage(req *dto.SendMessageDto, from int64, to string) (int, error)
	FetchMessage(from int64, to string) ([]*dto.FetchMessageDto, int, error)
	SearchMessage(userId int64, message string) ([]*entity.Message, int, error)
}

type service struct {
	userChatRepository    repository.UserChatRepository
	messageChatRepository repository.MessageRepository
	kafkaProducer         kafka.Producer
	validator             validator.Validator
}

func NewChatService(
	userChatRepository repository.UserChatRepository,
	messageRepository repository.MessageRepository,
	kafkaProducer kafka.Producer,
	validator validator.Validator,
) Service {
	return &service{
		userChatRepository:    userChatRepository,
		messageChatRepository: messageRepository,
		kafkaProducer:         kafkaProducer,
		validator:             validator,
	}
}

func (s *service) SendMessage(req *dto.SendMessageDto, from int64, to string) (int, error) {
	if err := s.validator.Validate(req); err != nil {
		return customstatus.ErrBadRequest.Code, errors.New(customstatus.ErrBadRequest.Message)
	}

	toID, err := strconv.Atoi(to)
	if err != nil {
		return customstatus.ErrBadRequest.Code, errors.New(err.Error())
	}

	userChat, err := s.userChatRepository.GetUserChatByFromAndTo(from, int64(toID))
	if err != nil {
		userChatReq := &entity.UserChat{
			From: from,
			To:   int64(toID),
		}
		userChat, err = s.userChatRepository.CreateUserChat(userChatReq)
		if err != nil {
			return customstatus.ErrInternalServerError.Code, errors.New(err.Error())
		}
	}

	jsonMessage, err := json.Marshal(dto.KafkaSendMessage{
		UserChatId:  userChat.ID,
		MessageText: req.Message,
	})
	if err != nil {
		return customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}

	kafkaMessages := &sarama.ProducerMessage{
		Topic: config.GetEnv("CHAT_TOPIC"),
		Value: sarama.StringEncoder(jsonMessage), // Use StringEncoder for JSON data
	}
	err = s.kafkaProducer.SendMessage(kafkaMessages)
	if err != nil {
		return customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}

	return customstatus.StatusCreated.Code, nil
}

func (s *service) FetchMessage(from int64, to string) ([]*dto.FetchMessageDto, int, error) {
	toID, err := strconv.Atoi(to)
	if err != nil {
		return nil, customstatus.ErrBadRequest.Code, errors.New(err.Error())
	}

	userChat, err := s.userChatRepository.GetUserChatByORFromTo(from, int64(toID))
	if err != nil {
		return nil, customstatus.StatusOk.Code, nil
	}

	query := fmt.Sprintf("user_chat_id.eq.%d", userChat[0].ID)
	if len(userChat) > 1 {
		query = fmt.Sprintf("user_chat_id.eq.%d,user_chat_id.eq.%d", userChat[0].ID, userChat[1].ID)
	}

	data, err := helpers.GetMessageFromSupabase(query)
	if err != nil {
		return nil, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}

	return data, customstatus.StatusOk.Code, nil
}

func (s *service) SearchMessage(userId int64, message string) ([]*entity.Message, int, error) {
	messageChat := s.messageChatRepository.SearchMessage(userId, message)
	fmt.Println(messageChat)
	if messageChat == nil {
		return nil, customstatus.StatusOk.Code, nil
	}

	return messageChat, customstatus.StatusOk.Code, nil
}
