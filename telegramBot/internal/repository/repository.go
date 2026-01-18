package repository

import (
	"controller/pkg/models"
	"database/sql"
)

type UserRepo interface {
	GetUserById(userId int64) (*models.User, error)
	CreateUser(userId int64, username string) error
	SetSubscribed(userId int64, subscribeStatus int) (bool, error)
}

type Repository struct {
	UserRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepo: NewUserPostgres(db),
	}
}
