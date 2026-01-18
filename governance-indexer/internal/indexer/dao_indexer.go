package indexer

import (
	"bytes"
	"controller/pkg/models"
	"controller/pkg/service"
	"encoding/json"
	"fmt"
	"governance-indexer/internal/repository"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type DAOIndexer struct {
	repo    *repository.Repository
	rwKafka *service.ReaderWriterService
}

func NewDAOIndexer(repo *repository.Repository, rwKafka *service.ReaderWriterService) *DAOIndexer {
	return &DAOIndexer{repo: repo, rwKafka: rwKafka}
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

var querySpaces = `
{
	spaces(
		first: %d,
		orderBy: "created",
		orderDirection: desc
	) {
		id
		name
		about
		network
		symbol
		created
		admins
		members
		filters {
			minScore
			onlyMembers
		}
		strategies {
			name
		}
	}
}`

var queryVotes = `
{
	votes (
		first: %d
		skip: %d
		where: {
			proposal: "%s"
		}
		orderBy: "created",
		orderDirection: desc
	) {
		id
		voter
		created
		choice
		vp
		vp_state
	}
}`

// CreatedResponse и DataResponse Структуры для получения ответа
type CreatedResponse struct {
	Data json.RawMessage `json:"data"`
}

// Структуры для получения ответа
type ProposalsResponse struct {
	Proposals []models.Proposals `json:"proposals"`
}

type SpaceResponse struct {
	Space []models.Space `json:"spaces"`
}

type VotesResponse struct {
	Votes []models.Votes `json:"votes"`
}

func (d *DAOIndexer) MainIndex(numberRecords int, typeQuery string) error {

	var query string

	switch typeQuery {
	case "proposals":
		{
			log.Println("Proposals indexer")
			query = fmt.Sprintf(queryProposals, numberRecords)
		}
	case "spaces":
		{
			log.Println("Spaces indexer")
			query = fmt.Sprintf(querySpaces, numberRecords)
		}
	case "votes":
		{
			// FIXME:Тесть. Пока не понимаю как это подключить. Нужно Proposals для работы
			log.Println("Votes indexer")
			err := d.RequestVotes(1000, "QmPvbwguLfcVryzBRrbY4Pb9bCtxURagdv1XjhtFLf3wHj")
			if err != nil {
				return err
			}
			return nil
		}
	case "getAllSpaces":
		{
			log.Println("getAllSpaces indexer")
			err := d.RequestSpaces(1000, 20)
			if err != nil {
				return err
			}
			return nil
		}
	}

	err := d.Request(query)
	if err != nil {
		return err
	}
	return nil
}

var querySpacesNum = `
{
	spaces(
		first: %d,
		skip: %d,
		orderDirection: desc
	) {
		id
		name
		about
		network
		symbol
		created
		admins
		members
		filters {
			minScore
			onlyMembers
		}
		strategies {
			name
		}
	}
}`

func (d *DAOIndexer) RequestSpaces(batchSize int, sleepTime time.Duration) error {
	var lenGet = batchSize
	var batchN = 0
	var query = ""
	var resLen = 0

	for {
		skip := batchSize * batchN
		query = fmt.Sprintf(querySpacesNum, batchSize, skip)

		log.Println("Request graph query:", query)
		// Переводим в json формат
		jsonData, err := json.Marshal(map[string]string{
			"query": query,
		})
		if err != nil {
			log.Println("JSON marshal customErrors:", err)
			return err
		}

		// Отправляем запрос и получаем ответ
		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("HTTP request customErrors:", err)
			return err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		// получаем тело ответа
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Read body customErrors:", err)
			return err
		}

		// парсим в json
		var response CreatedResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Println("JSON unmarshal customErrors:", err)
			return err
		}

		var spaces SpaceResponse
		if err := json.Unmarshal(response.Data, &spaces); err == nil {
			err := d.repo.SpaceRepo.AddHistory(spaces.Space)
			if err != nil {
				return err
			}
			fmt.Println("Len Array", spaces.Space[0])
		}

		batchN += 1
		lenGet = len(spaces.Space)
		resLen += lenGet
		fmt.Println("Result len ", resLen)

		if batchSize != lenGet {
			break
		}
		time.Sleep(sleepTime * time.Millisecond)
	}
	fmt.Println("End")
	return nil
}

func (d *DAOIndexer) RequestVotes(batchSize int, proposals string) error {
	var lenGetVotes = batchSize
	var batchN = 0
	var query = ""
	var resLen = 0

	for {
		skip := batchSize * batchN
		query = fmt.Sprintf(queryVotes, batchSize, skip, proposals)

		log.Println("Request graph query:", query)
		// Переводим в json формат
		jsonData, err := json.Marshal(map[string]string{
			"query": query,
		})
		if err != nil {
			log.Println("JSON marshal customErrors:", err)
			return err
		}

		// Отправляем запрос и получаем ответ
		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("HTTP request customErrors:", err)
			return err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		// получаем тело ответа
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Read body customErrors:", err)
			return err
		}

		//fmt.Println(string(body))

		// парсим в json
		var response CreatedResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Println("JSON unmarshal customErrors:", err)
			return err
		}

		//fmt.Println("String", string(response.Data))
		//fmt.Println("Len", len(response.Data))

		var votes VotesResponse
		if err := json.Unmarshal(response.Data, &votes); err == nil {

			fmt.Println("Len Votes", votes.Votes[0])
		}

		fmt.Println("lenGetVotes ", lenGetVotes, "==", batchSize)
		fmt.Println("batchN ", batchN)

		batchN += 1
		lenGetVotes = len(votes.Votes)
		resLen += lenGetVotes
		fmt.Println("Result len ", resLen)
		break
		//if batchSize != lenGetVotes {
		//	break
		//}
		//time.Sleep(20 * time.Millisecond)
	}
	fmt.Println("End")
	return nil
}

// IndexProposal получает записи proposal и сохраняет в БД
func (d *DAOIndexer) Request(graphQuery string) error {

	log.Println("Request graph query:", graphQuery)
	// Переводим в json формат
	jsonData, err := json.Marshal(map[string]string{
		"query": graphQuery,
	})
	if err != nil {
		log.Println("JSON marshal customErrors:", err)
		return err
	}

	// Отправляем запрос и получаем ответ
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("HTTP request customErrors:", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// получаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read body customErrors:", err)
		return err
	}

	// парсим в json
	var response CreatedResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("JSON unmarshal customErrors:", err)
		return err
	}

	//fmt.Println("response.Data", string(response.Data))

	// если это Proposals
	var proposals ProposalsResponse
	if err := json.Unmarshal(response.Data, &proposals); err == nil {
		// смотрим каких записей нет
		log.Println("Start Proposals")
		err := d.ProposalsProcessing(proposals.Proposals)
		if err != nil {
			return err
		}
	}

	//fmt.Println(response)

	// если это Proposals
	var spaceResponse SpaceResponse
	if err := json.Unmarshal(response.Data, &spaceResponse); err == nil {
		// смотрим каких записей нет
		log.Println("Start Space")
		err := d.SpaceProcessing(spaceResponse.Space)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DAOIndexer) ProposalsProcessing(proposals []models.Proposals) error {
	missing, err := d.repo.ProposalRepo.FindMissing(proposals)
	if err != nil {
		log.Println("Error finding missing proposals:", err)
		return err
	}

	lenMissing := len(missing)
	if lenMissing == 0 {
		log.Println("No missing proposals found")
		return nil
	}

	// TODO: здесь должен быть outbox-паттерн. Нужно написать часть проверяющую доставку сообщений и /
	// TODO: отправляет запрос заново, если подтверждения нет
	ids := make([]string, 0, lenMissing)
	for _, proposal := range missing {
		ids = append(ids, proposal.ID)
	}

	eventData := models.NewData{
		TableName: "proposals",
		IDs:       ids,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		log.Println(fmt.Sprintf("Marshal customErrors: %v", err))
	}
	err = d.rwKafka.WriteMessage(
		kafka.Message{
			Value: data,
		})
	if err != nil {
		log.Println(fmt.Sprintf("Write messages customErrors: %v", err))
	}

	log.Println("Proposals write:", len(missing))

	// если есть записи, то сохраняем в БД
	if len(missing) > 0 {
		if err := d.repo.ProposalRepo.AddProposal(missing); err != nil {
			log.Println("Repository customErrors:", err)
			return err
		}
	}
	return nil
}

func (d *DAOIndexer) SpaceProcessing(spaces []models.Space) error {
	missing, err := d.repo.SpaceRepo.FindMissing(spaces)
	if err != nil {
		log.Println("Error finding missing proposals:", err)
		return err
	}

	lenMissing := len(missing)
	if lenMissing == 0 {
		log.Println("No missing proposals found")
		return nil
	}

	// TODO: здесь должен быть outbox-паттерн. Нужно написать часть проверяющую доставку сообщений и /
	// TODO: отправляет запрос заново, если подтверждения нет
	ids := make([]string, 0, lenMissing)
	for _, space := range missing {
		ids = append(ids, space.ID)
	}

	eventData := models.NewData{
		TableName: "spaces",
		IDs:       ids,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		log.Println(fmt.Sprintf("Marshal customErrors: %v", err))
	}
	err = d.rwKafka.WriteMessage(
		kafka.Message{
			Value: data,
		})
	if err != nil {
		log.Println(fmt.Sprintf("Write messages customErrors: %v", err))
	}

	log.Println("Spaces write:", len(missing))

	// если есть записи, то сохраняем в БД
	if len(missing) > 0 {
		if err := d.repo.SpaceRepo.AddNew(missing); err != nil {
			log.Println("Repository customErrors:", err)
			return err
		}
	}
	return nil
}
