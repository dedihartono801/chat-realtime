package repository

import (
	"database/sql"

	"github.com/dedihartono801/chat-realtime/internal/entity"
)

type UserChatRepository interface {
	GetUserChatByFromAndTo(from int64, to int64) (*entity.UserChat, error)
	GetUserChatByORFromTo(from int64, to int64) ([]*entity.UserChat, error)
	CreateUserChat(userChat *entity.UserChat) (*entity.UserChat, error)
}

type userChatRepository struct {
	db *sql.DB
}

func NewUserChatRepository(db *sql.DB) UserChatRepository {
	return &userChatRepository{db}
}

func (r *userChatRepository) GetUserChatByFromAndTo(from int64, to int64) (*entity.UserChat, error) {
	query := `SELECT id, "from", "to" FROM user_chat WHERE "from" = $1 and "to" = $2`
	dtUserChat := &entity.UserChat{}

	err := r.db.QueryRow(query, from, to).Scan(&dtUserChat.ID, &dtUserChat.From, &dtUserChat.To)
	if err != nil {
		return nil, err
	}

	return dtUserChat, nil
}

func (r *userChatRepository) GetUserChatByORFromTo(from int64, to int64) ([]*entity.UserChat, error) {
	query := `SELECT id, "from", "to" FROM user_chat WHERE ("from" = $1 AND "to" = $2) OR ("from" = $2 AND "to" = $1)`
	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userChats []*entity.UserChat

	for rows.Next() {
		dtUserChat := &entity.UserChat{}
		var from, to sql.NullInt64 // Use NullInt64 to handle possible NULL values in database columns

		if err := rows.Scan(&dtUserChat.ID, &from, &to); err != nil {
			return nil, err
		}

		// Set values to UserChat if they are not NULL
		if from.Valid {
			dtUserChat.From = from.Int64
		}

		if to.Valid {
			dtUserChat.To = to.Int64
		}

		userChats = append(userChats, dtUserChat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return userChats, nil
}

func (r *userChatRepository) CreateUserChat(userChat *entity.UserChat) (*entity.UserChat, error) {
	query := `
        INSERT INTO user_chat ("from", "to")
        VALUES ($1, $2)
		RETURNING id, "from", "to", created_at
    `

	var userChatDt entity.UserChat
	err := r.db.QueryRow(query, userChat.From, userChat.To).Scan(&userChatDt.ID, &userChatDt.From, &userChatDt.To, &userChatDt.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &userChatDt, nil
}
