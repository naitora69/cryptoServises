package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	customError "github.com/danila-kuryakin/cryptoServises/controller/errors"
	"github.com/danila-kuryakin/cryptoServises/controller/internal/repository"
	"github.com/danila-kuryakin/cryptoServises/controller/pkg/models"
	"github.com/danila-kuryakin/cryptoServises/controller/pkg/service"

	"github.com/segmentio/kafka-go"
)

type DaoService struct {
	repo         *repository.Repository
	indexerKafka *service.ReaderWriterService
	botKafka     *service.ReaderWriterService
}

func NewDaoService(repo *repository.Repository, indexerKafka *service.ReaderWriterService, botKafka *service.ReaderWriterService) *DaoService {
	return &DaoService{repo: repo, indexerKafka: indexerKafka, botKafka: botKafka}
}

func (d *DaoService) Processing() error {
	message, err := d.indexerKafka.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
	}
	log.Println(message)

	var eventData models.NewData
	if err := json.Unmarshal(message.Value, &eventData); err != nil {
		log.Println("bad json:", err)
	}

	if eventData.TableName == "proposals" {
		if err := d.ProcessingProposals(eventData); err != nil {
			log.Println("proposal processing failed: ", err.Error())
			return err
		}
	}

	return nil
}

func compareIDs(a []models.ProposalEvent, b []string) (same, onlyA, onlyB []string) {
	m := make(map[string]bool, len(a))

	for _, id := range a {
		m[id.ID] = true
	}

	for _, id := range b {
		if m[id] {
			same = append(same, id)
			delete(m, id) // чтобы не попал в onlyA
		} else {
			onlyB = append(onlyB, id)
		}
	}

	for id := range m {
		onlyA = append(onlyA, id)
	}

	return
}

func (d *DaoService) ProcessingProposals(proposalsData models.NewData) error {
	// получаем события из бд
	proposalsNew, err := d.repo.ReadProposalsEvents()
	if err != nil {
		log.Println("Error reading proposals: ", err)
		return err
	}

	// Проверяем не соответствие id для логов
	_, onlyA, onlyB := compareIDs(proposalsNew, proposalsData.IDs)
	if len(onlyA) > 0 {
		log.Println("Error: Arrays A do not match.", onlyA)
	}
	if len(onlyB) > 0 {
		log.Println("Error: Arrays B do not match.", onlyB)
	}

	// добавляем события в планировщик
	if err := d.repo.AddEventScheduler(proposalsNew); err != nil {
		log.Println("Error adding event scheduler: ", err)
		return err
	}

	// Подтверждаем прочтение сообщения
	if err := d.repo.ProposalDeliverySuccessful(proposalsNew); err != nil {
		log.Println("delivery customErrors:", err)
		return err
	}
	return nil
}

func (d *DaoService) MessageController() error {
	events, err := d.repo.GetCurrentEvents(3)
	if err != nil {
		if errors.Is(err, customError.ErrDataNotFound) {
			log.Println("No data found.")
			return nil
		}
		log.Println("Error getting current events: ", err)
		return err
	}
	subscriptions, err := d.repo.GetUserSubscriptions()
	if err != nil {
		return err
	}
	for _, event := range events {
		if event.EventTime.Time.Before(time.Now()) {
			fmt.Println(event)
			msgData := models.CurrentProposalEvent{Users: subscriptions, CurrentEvent: event}
			data, err := json.Marshal(msgData)
			if err != nil {
				log.Println(fmt.Sprintf("Marshal customErrors: %v", err))
			}

			err = d.botKafka.WriteMessage(kafka.Message{
				Value: data,
			})
			if err != nil {
				log.Println(fmt.Sprintf("WriteMessage err: %v", err))
				return err
			}
			if err := d.repo.EventDeliverySuccessful([]models.CurrentEvent{event}); err != nil {
				return err
			}
		}
	}

	return nil
}
