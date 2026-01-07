package service

import (
	"controller/internal/config"
	"controller/internal/repository"
	"controller/pkg/service"

	"github.com/segmentio/kafka-go"
)

type Dao interface {
}

type ReaderWriterKafka interface {
	ReadMessage() (*kafka.Message, error)
	WriteMessage(message kafka.Message) error
}

type Service struct {
	Dao
	ReaderWriterKafka
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		Dao:               NewDaoService(repo),
		ReaderWriterKafka: service.NewReaderWriterService(cfg.Kafka.Address, cfg.Kafka.Port, config.DaoIndexerTopic, config.DaoIndexerGroup),
	}
}
