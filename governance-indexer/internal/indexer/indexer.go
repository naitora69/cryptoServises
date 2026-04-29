package indexer

import (
	"time"

	"github.com/danila-kuryakin/cryptoServises/controller/pkg/models"
	"github.com/danila-kuryakin/cryptoServises/controller/pkg/service"
	"github.com/danila-kuryakin/cryptoServises/governance-indexer/internal/config"
	"github.com/danila-kuryakin/cryptoServises/governance-indexer/internal/repository"
)

type DAOIndexerInterface interface {
	MainIndex(numberRecords int, typeQuery string) error
	Request(graphQuery string) error

	RequestSpaces(batchSize int, sleepTime time.Duration) error
	//RequestProposals(batchSize int) error
	RequestVotes(batchSize int, proposals string) error

	ProposalsProcessing(proposals []models.Proposals) error
	SpaceProcessing(space []models.Space) error
}

type Indexer struct {
	DAOIndexerInterface
}

func NewIndexer(repo *repository.Repository, cfg *config.Config) *Indexer {
	rwKafka := service.NewReaderWriterService(cfg.Kafka.Address, cfg.Kafka.Port, config.DaoIndexerTopic, config.DaoIndexerGroup)
	return &Indexer{
		DAOIndexerInterface: NewDAOIndexer(repo, rwKafka),
	}
}
