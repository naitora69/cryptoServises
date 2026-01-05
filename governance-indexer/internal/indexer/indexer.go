package indexer

import (
	"governance-indexer/internal/repository"

	"github.com/segmentio/kafka-go"
)

type ProposalIndexerInterface interface {
	IndexProposal(numberRecords int) error
}

type Indexer struct {
	ProposalIndexerInterface
}

func NewIndexer(repo *repository.Repository, writer *kafka.Writer) *Indexer {
	return &Indexer{
		ProposalIndexerInterface: NewProposalIndexer(repo, writer),
	}
}
