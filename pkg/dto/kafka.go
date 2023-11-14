package dto

type KafkaSendMessage struct {
	UserChatId  int64  `json:"user_chat_id"`
	MessageText string `json:"message_text"`
}
