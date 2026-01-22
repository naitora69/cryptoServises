package timer

import (
	"governance-indexer/internal/config"
	"governance-indexer/internal/indexer"
	"log"
	"time"
)

type ProposalTimer struct {
	index  *indexer.Indexer
	config *config.Config
}

func NewProposalTimer(index *indexer.Indexer, config *config.Config) *ProposalTimer {
	return &ProposalTimer{index: index, config: config}
}

func (p ProposalTimer) StartProposal() {
	for {
		err := p.index.MainIndex(p.config.Proposal.NumberRecords, "proposals")
		if err != nil {
			log.Println("Error indexing proposal records", err)
			return
		}
		var durationMinutes = time.Duration(p.config.Proposal.TimeRequest)
		time.Sleep(durationMinutes * time.Second)
	}
}

func (p ProposalTimer) StartSpace(saveAllSpaces bool) {

	// получаем все spaces если нужно
	if saveAllSpaces {
		err := p.index.MainIndex(p.config.Proposal.NumberRecords, "getAllSpaces")
		if err != nil {
			log.Println("Error indexing getAllSpaces records", err)
			return
		}
	}

	// получаем новые события spaces
	for {
		err := p.index.MainIndex(p.config.Proposal.NumberRecords, "spaces")
		if err != nil {
			log.Println("Error indexing space records", err)
			return
		}
		var durationMinutes = time.Duration(p.config.Proposal.TimeRequest)
		time.Sleep(durationMinutes * time.Second)
	}
}

func (p ProposalTimer) StartVotes() {

	err := p.index.MainIndex(p.config.Proposal.NumberRecords, "votes")
	if err != nil {
		log.Println("Error indexing votes records", err)
		return
	}
}
