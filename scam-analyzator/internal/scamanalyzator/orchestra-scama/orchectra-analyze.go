package orchestrascama

import (
	"log"
	ageanalyze "scam-analyzator-service/internal/scamanalyzator/age-analyze"
	countholders "scam-analyzator-service/internal/scamanalyzator/count-holders"
	"scam-analyzator-service/internal/scamanalyzator/fetch"
)

type ScamsFlags struct {
	TokenAgeScam          bool
	TokenHoldersCountScam bool
}

func NewScamsFlags(tas bool, thcs bool) ScamsFlags {
	return ScamsFlags{
		TokenAgeScam:          tas,
		TokenHoldersCountScam: thcs,
	}
}

func FinalRes(spaceID string) (ScamsFlags, error) {
	tokenIds, errors := fetch.GetTokenId(spaceID)
	if errors != nil {
		log.Println("Не найден токены для данного пр-ва")
		return ScamsFlags{}, errors[0]
	}
	ageFlag := ageanalyze.GetTokenAgeFlag(tokenIds)
	countHoldersFlag, err := countholders.GetTokenHoldersCountFlag(tokenIds)
	if err != nil {
		log.Println(err)
	}
	return NewScamsFlags(ageFlag, countHoldersFlag), nil
}
