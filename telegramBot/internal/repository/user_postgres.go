package repository

import (
	"controller/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	error2 "telegramBot/internal/error"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (p *UserPostgres) GetUserById(userId int64) (*models.User, error) {

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, usersTable)

	var user models.User
	err := p.db.QueryRow(query, userId).Scan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error2.ErrorUserNotFound
		}
		return nil, errors.New(fmt.Sprint("Error getting user by id ", userId))
	}
	return &user, nil
}

func (p *UserPostgres) CreateUser(userId int64, username string) error {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, dao_subscribed) VALUES ($1, $2)`, usersTable)
	_, err := p.db.Exec(query, userId, 0)

	if err != nil {
		return errors.New(fmt.Sprint("Error creating user ", userId))
	}
	return nil
}
func (p *UserPostgres) SetSubscribed(userId int64, subscribeStatus int) (bool, error) {
	query := fmt.Sprintf(`
        UPDATE %s SET dao_subscribed = $1
        WHERE user_id = $2`, usersTable)

	result, err := p.db.Exec(query, subscribeStatus, userId)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, error2.ErrorUserNotFound
	}
	return true, nil
}
