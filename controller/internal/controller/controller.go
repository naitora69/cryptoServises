package controller

import (
	"log"
	"time"

	"github.com/danila-kuryakin/cryptoServises/controller/internal/service"
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
	for {
		if err := c.service.Processing(); err != nil {
			log.Println(err)
		}
	}
}

func (c *Controller) InitMessageController() {
	for {
		if err := c.service.MessageController(); err != nil {
			return
		}

		time.Sleep(1 * time.Second)
	}
}
