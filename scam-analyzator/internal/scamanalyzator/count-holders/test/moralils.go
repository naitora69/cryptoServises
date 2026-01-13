package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"scam-analyzator-service/internal/config"
	"scam-analyzator-service/internal/scamanalyzator"
)

func main() {

	testspace := "stgdao.eth"
	tokenData, err := scamanalyzator.GetTokenId(testspace)
	if err != nil {
		log.Fatal(err)
	}
	_, apiKey := config.GetApiKey()

	chain := scamanalyzator.GetChainName(tokenData.Network[0])
	for _, v := range tokenData.TokensAddress {
		fmt.Printf("Проверяем адрес: %s в сети %s...\n", v, chain)

		url := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/erc20/%s/owners?chain=%s", v, chain)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("X-API-Key", apiKey)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("HttpClient error: ", err)
		}

		body, _ := io.ReadAll(res.Body)
		res.Body.Close()

		if !bytes.Contains(body, []byte(`"result":[]`)) {
			fmt.Println("УРА! Нашли живой токен.")
			fmt.Println(string(body))
		} else {
			fmt.Println("Пусто, проверяем следующий...")
		}
	}

}
