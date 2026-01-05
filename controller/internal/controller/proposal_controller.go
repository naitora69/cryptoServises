package controller

import (
	"controller/internal/service"
	"net/http"
)

type ProposalController struct {
	service *service.Service
}

func NewProposalController(service *service.Service) *ProposalController {
	return &ProposalController{service: service}
}

func (p *ProposalController) Create(w http.ResponseWriter, r *http.Request) {
	return
}
