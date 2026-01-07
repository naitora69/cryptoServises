package scamanalyzator

import (
	"log"
	"scam-analyzator-service/internal/config"
	"time"
)

func FinalRes(tokenStruct TokenIdAnswer) bool {
	isScamFound := false
	ethKey := config.GetApiKey()

	for _, addr := range tokenStruct.tokensAddress {

		res, err := GetTokenAge(ethKey, addr, tokenStruct.network)
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
