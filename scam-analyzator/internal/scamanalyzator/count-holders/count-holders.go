package countholders

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scam-analyzator-service/internal/config"
	"scam-analyzator-service/internal/scamanalyzator/fetch"
)

type TokenHoldersResponce struct {
	TotalHolders int `json:"totalHolders"`
}

// TODO продумать над тем как возвращать флаг скама, касаемо ошибок впервую очередь
func GetTokenHoldersCountFlag(tokenData fetch.TokenIdAnswer) (bool, error) {
	var scamFound bool = false

	_, apiKey := config.GetApiKey()
	for _, v := range tokenData.Tokens {
		var tokenResponce TokenHoldersResponce
		url := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/erc20/%s/holders?chain=%s", v.Address, v.Chain)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("X-API-Key", apiKey)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("HttpClient error: ", err)
			return false, err
		}

		body, _ := io.ReadAll(res.Body)
		res.Body.Close()
		err = json.Unmarshal(body, &tokenResponce)
		if err != nil {
			log.Println("Unmarshall error: ", err)
			return false, err
		}
		if tokenResponce.TotalHolders < 50 {
			log.Println(v.Address, " ", v.Chain, " < 50 holders")
			scamFound = true
		}

	}
	return scamFound, nil
}
