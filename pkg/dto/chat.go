package dto

type SendMessageDto struct {
	Message string `json:"message" validate:"required"`
}

type FetchMessageDto struct {
	ID          int64  `json:"id" validate:"required"`
	UserChatID  int64  `json:"user_chat_id" validate:"required"`
	MessageText string `json:"message_text" validate:"required"`
}

type SaveMessageDto struct {
	UserChatID  int64  `json:"user_chat_id" validate:"required"`
	MessageText string `json:"message_text" validate:"required"`
}
