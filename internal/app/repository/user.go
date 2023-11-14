package repository

import (
	"database/sql"

	"github.com/dedihartono801/chat-realtime/internal/entity"
)

type UserRepository interface {
	GetUserByUsername(username string) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByUsername(username string) (*entity.User, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"
	user := &entity.User{}

	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) CreateUser(user *entity.User) (*entity.User, error) {
	query := `
        INSERT INTO users (name, username, password)
        VALUES ($1, $2, $3)
        RETURNING id, name, username, password, created_at
    `

	var insertedUser entity.User
	err := r.db.QueryRow(query, user.Name, user.Username, user.Password).Scan(&insertedUser.ID, &insertedUser.Name, &insertedUser.Username, &insertedUser.Password, &insertedUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}
