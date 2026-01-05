package repository

import (
	"database/sql"
)

type ProposalRepo interface {
}

type Repository struct {
	ProposalRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ProposalRepo: NewProposalPostgres(db),
	}
}
