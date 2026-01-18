package service

import (
	"controller/pkg/models"
	pkgService "controller/pkg/service"
	"telegramBot/internal/config"
	"telegramBot/internal/repository"
)

type Dao interface {
	NewUser(userId int64, username string) error
	Subscribed(userId int64) (bool, error)
	Unsubscribed(userId int64) (bool, error)
	KafkaListen() (models.CurrentProposalEvent, error)
}

type Service struct {
	Dao
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	controllerKafka := pkgService.NewReaderWriterService(
		cfg.Kafka.Address,
		cfg.Kafka.Port,
		config.DaoControllerBotTopic,
		config.DaoControllerBotGroup)
	return &Service{
		Dao: NewDaoService(repo, controllerKafka),
	}
}
