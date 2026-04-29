package service

import (
	"github.com/danila-kuryakin/cryptoServises/controller/pkg/models"
	pkgService "github.com/danila-kuryakin/cryptoServises/controller/pkg/service"
	"github.com/danila-kuryakin/cryptoServises/telegramBot/internal/config"
	"github.com/danila-kuryakin/cryptoServises/telegramBot/internal/repository"
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
