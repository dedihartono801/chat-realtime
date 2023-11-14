package repository

import (
	"database/sql"

	"github.com/dedihartono801/chat-realtime/internal/entity"
)

type MessageRepository interface {
	CreateMessage(message *entity.Message) error
	SearchMessage(userId int64, message string) []*entity.Message
}

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) CreateMessage(message *entity.Message) error {
	query := `
        INSERT INTO messages (user_chat_id, message_text)
        VALUES ($1, $2)
    `

	_, err := r.db.Exec(query, message.UserChatID, message.MessageText)
	if err != nil {
		return err
	}

	return nil
}

func (r *messageRepository) SearchMessage(userId int64, message string) []*entity.Message {
	query := `
        select messages.id, messages.user_chat_id, messages.message_text, messages.created_at from messages 
		join user_chat on messages.user_chat_id = user_chat.id
		WHERE 
		(user_chat.from = $1 OR user_chat.to = $1) 
		AND
		messages.message_text LIKE '%' || $2 || '%'
    `
	rows, err := r.db.Query(query, userId, message)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var messageChat []*entity.Message

	for rows.Next() {
		dtMessage := &entity.Message{}

		if err := rows.Scan(&dtMessage.ID, &dtMessage.UserChatID, &dtMessage.MessageText, &dtMessage.CreatedAt); err != nil {
			return nil
		}

		messageChat = append(messageChat, dtMessage)
	}

	if err = rows.Err(); err != nil {
		return nil
	}
	return messageChat
}
