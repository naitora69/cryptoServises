package timer

import (
	"governance-indexer/internal/config"
	"governance-indexer/internal/indexer"
)

type ProposalTimerInterface interface {
	StartProposal()
	StartSpace(saveAllSpaces bool)
	StartVotes()
}

type Timer struct {
	ProposalTimerInterface
}

func NewTimer(index *indexer.Indexer, cfg *config.Config) *Timer {
	return &Timer{
		ProposalTimerInterface: NewProposalTimer(index, cfg),
	}
}
