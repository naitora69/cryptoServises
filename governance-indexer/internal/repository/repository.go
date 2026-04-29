package repository

import (
	"database/sql"

	"github.com/danila-kuryakin/cryptoServises/controller/pkg/models"
)

type ProposalRepo interface {
	AddProposal(proposals []models.Proposals) error
	FindMissing(proposals []models.Proposals) ([]models.Proposals, error)
}

type SpaceRepo interface {
	Add(space []models.Space, eventType string) error
	AddHistory(space []models.Space) error
	AddNew(space []models.Space) error
	FindMissing(spaces []models.Space) ([]models.Space, error)
}

type Repository struct {
	ProposalRepo
	SpaceRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ProposalRepo: NewProposalPostgres(db),
		SpaceRepo:    NewSpacePostgres(db),
	}
}
