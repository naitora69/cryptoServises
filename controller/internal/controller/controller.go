package controller

import (
	"controller/internal/service"
	"controller/pkg/models"
	"encoding/json"
	"log"
)

type Controller struct {
	service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) InitListener() {
	i := 1
	for {
		m, err := c.service.ReaderWriterKafka.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		log.Println(m)

		var proposals_in []models.Proposals
		if err := json.Unmarshal(m.Value, &proposals_in); err != nil {
			log.Println("bad json:", err)
			continue
		}

		log.Printf("  %d) received order: %#v\n", i, proposals_in)
		i++
	}
}
