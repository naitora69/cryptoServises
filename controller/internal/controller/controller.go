package controller

import (
	"context"
	"controller/internal/models"
	"controller/internal/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

type Proposal interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type Controller struct {
	proposal Proposal
	reader   *kafka.Reader
}

func NewController(service *service.Service, reader *kafka.Reader) *Controller {
	return &Controller{
		proposal: NewProposalController(service),
		reader:   reader,
	}
}

func (c *Controller) InitListener() {
	i := 1
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		var proposals_in []models.Proposals
		if err := json.Unmarshal(m.Value, &proposals_in); err != nil {
			log.Println("bad json:", err)
			continue
		}

		log.Printf("  %d) received order: %#v\n", i, proposals_in)
		i++
	}
}
