package ageanalyze

import (
	"log"
	"scam-analyzator-service/internal/config"
	"scam-analyzator-service/internal/scamanalyzator/fetch"
	"time"
)

func FinalRes(tokenStruct fetch.TokenIdAnswer) bool {
	isScamFound := false

	_, moralisKey := config.GetApiKey()

	for _, addr := range tokenStruct.Tokens {
		res, err := GetTokenAge(moralisKey, addr.Address, addr.Chain)
		if err != nil {
			log.Printf("Age Check error for %s (%s): %v", addr.Address, addr.Chain, err)
			continue
		}

		analysis := AnalyzeAge(res)
		if analysis.IsNewToken {
			log.Printf("ВНИМАНИЕ: Токен %s в сети %s ОПАСЕН. Возраст: %.2f ч.",
				addr.Address, addr.Chain, analysis.AgeHours)
			isScamFound = true
		} else {
			log.Printf("Токен %s (%s) прошел проверку. Возраст: %.2f дней",
				addr.Address, addr.Chain, analysis.AgeHours/24)
		}

		time.Sleep(200 * time.Millisecond)
	}
	return isScamFound
}
