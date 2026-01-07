package indexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"governance-indexer/internal/repository"
	"governance-indexer/pkg/models"
	"governance-indexer/pkg/service"
	"io"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

type ProposalIndexer struct {
	repo    *repository.Repository
	rwKafka *service.ReaderWriterService
}

func NewProposalIndexer(repo *repository.Repository, rwKafka *service.ReaderWriterService) *ProposalIndexer {
	return &ProposalIndexer{repo: repo, rwKafka: rwKafka}
}

var endpoint = "https://hub.snapshot.org/graphql"

// GraphQL-запрос
var queryProposals = `
{
	proposals (
		first: %d,
		skip: 0,
		orderDirection: desc
	) {
		id
		title
		author
		created
		start
		end
		state
		snapshot
		choices
		space {
			id
			name
		}
	}
}`

// CreatedResponse и DataResponse Структуры для парсинга ответа
type CreatedResponse struct {
	Data DataResponse `json:"data"`
}

type DataResponse struct {
	Proposals []models.Proposals `json:"proposals"`
}

// IndexProposal получает записи proposal и сохраняет в БД
func (p *ProposalIndexer) IndexProposal(numberRecords int) error {

	// Дописываем запрос. Добавляем количество получаемых записей
	query := fmt.Sprintf(queryProposals, numberRecords)

	// Переводим в json формат
	jsonData, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		log.Println("JSON marshal error:", err)
		return err
	}

	// Отправляем запрос и получаем ответ
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("HTTP request error:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read body error:", err)
		return err
	}

	// парсим в json
	var result CreatedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("JSON unmarshal error:", err)
		return err
	}

	// смотрим каких записей нет
	missing, err := p.repo.FindMissing(result.Data.Proposals)
	if err != nil {
		log.Println("Error finding missing proposals:", err)
		return err
	}

	if len(missing) == 0 {
		log.Println("No missing proposals found")
		return nil
	}
	// TODO: здесь должен быть outbox-паттерн. Пока использую только kafka

	fmt.Println(missing)
	data, err := json.Marshal(missing)
	if err != nil {
		log.Fatal(fmt.Sprintf("Marshal error: %v", err))
	}
	err = p.rwKafka.WriteMessage(
		kafka.Message{
			Value: data,
		})
	if err != nil {
		log.Fatal(fmt.Sprintf("Write messages error: %v", err))
	}

	log.Println("Proposals writes:", len(missing))

	// если есть записи, то сохраняем в БД
	if len(missing) > 0 {
		if err := p.repo.AddProposal(missing); err != nil {
			log.Println("Repository error:", err)
			return err
		}
	}

	return nil
}
