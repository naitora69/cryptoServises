package service

import (
	"controller/internal/repository"
)

type DaoService struct {
	repo *repository.Repository
}

func NewDaoService(repo *repository.Repository) *DaoService {
	return &DaoService{repo: repo}
}
