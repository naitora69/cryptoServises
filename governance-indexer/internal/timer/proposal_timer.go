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
		err := p.index.IndexProposal(p.config.Proposal.NumberRecords)
		if err != nil {
			log.Println("Error indexing proposal records", err)
			return
		}

		var durationMinutes = time.Duration(p.config.Proposal.TimeRequest)

		time.Sleep(durationMinutes * time.Second)
	}
}
