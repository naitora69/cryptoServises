package indexer

import (
	"governance-indexer/internal/config"
	"governance-indexer/internal/repository"
	"governance-indexer/pkg/service"

	"github.com/segmentio/kafka-go"
)

type ProposalIndexerInterface interface {
	IndexProposal(numberRecords int) error
}

type ReaderWriterKafka interface {
	ReadMessage() (*kafka.Message, error)
	WriteMessage(message kafka.Message) error
}

type Indexer struct {
	ProposalIndexerInterface
}

func NewIndexer(repo *repository.Repository, cfg *config.Config) *Indexer {
	rwKafka := service.NewReaderWriterService(cfg.Kafka.Address, cfg.Kafka.Port, config.DaoIndexerTopic, config.DaoIndexerGroup)
	return &Indexer{
		ProposalIndexerInterface: NewProposalIndexer(repo, rwKafka),
	}
}
