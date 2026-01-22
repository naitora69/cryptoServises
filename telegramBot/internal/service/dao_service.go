package service

import (
	"controller/pkg/models"
	pkgService "controller/pkg/service"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	cusmomError "telegramBot/internal/error"
	"telegramBot/internal/repository"
)

type DaoService struct {
	repo            *repository.Repository
	controllerKafka *pkgService.ReaderWriterService
}

func NewDaoService(repo *repository.Repository, controllerKafka *pkgService.ReaderWriterService) *DaoService {
	return &DaoService{repo: repo, controllerKafka: controllerKafka}
}

func (s *DaoService) NewUser(userId int64, username string) error {
	_, err := s.repo.GetUserById(userId)
	if err != nil {
		if errors.Is(err, cusmomError.ErrorUserNotFound) {
			err = s.repo.CreateUser(userId, username)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (s DaoService) Subscribed(userId int64) (bool, error) {
	fmt.Println("Subscribed ", userId)
	_, err := s.repo.SetSubscribed(userId, 1)
	if err != nil {
		return false, err
	}
	log.Println("Subscribed", userId)
	return true, nil
}

func (s DaoService) Unsubscribed(userId int64) (bool, error) {
	_, err := s.repo.SetSubscribed(userId, 0)
	if err != nil {
		return false, err
	}
	log.Println("Unsubscribed", userId)
	return true, nil
}

func (s DaoService) KafkaListen() (models.CurrentProposalEvent, error) {

	message, err := s.controllerKafka.ReadMessage()
	if err != nil {
		log.Println("Error reading message", err)
		return models.CurrentProposalEvent{}, err
	}

	var eventData models.CurrentProposalEvent
	if err := json.Unmarshal(message.Value, &eventData); err != nil {
		log.Println("bad json:", err)
	}

	return eventData, nil
}
