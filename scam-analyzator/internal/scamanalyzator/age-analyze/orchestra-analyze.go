package ageanalyze

import (
	"log"
	"scam-analyzator-service/internal/config"
	"scam-analyzator-service/internal/scamanalyzator"
	"time"
)

func FinalRes(tokenStruct scamanalyzator.TokenIdAnswer) bool {
	isScamFound := false
	ethKey, _ := config.GetApiKey()

	for i, addr := range tokenStruct.TokensAddress {

		res, err := GetTokenAge(ethKey, addr, tokenStruct.Network[i])
		if err != nil {
			log.Println("EtherScan error for: ", addr, err)
			continue
		}
		analysis := AnalyzeAge(res)
		if analysis.IsNewToken {
			log.Printf("Token - %s слишком новый (%.2f ч.) - это подозрительно", addr, analysis.AgeHours)
			isScamFound = true
		}

		time.Sleep(400 * time.Millisecond)

	}
	return isScamFound
}
