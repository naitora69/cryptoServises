package service

import (
	"controller/internal/repository"
)

type Dao interface {
}

type Service struct {
	Dao
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Dao: NewDaoService(repo),
	}
}
