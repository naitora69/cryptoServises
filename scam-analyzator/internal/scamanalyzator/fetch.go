package scamanalyzator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var queryToken = `
{ 
  space(id: "%s") {
	network 
	strategies {
	  params 
	} 
  }
}`

type SnapshotResponce struct {
	Data struct {
		Space struct {
			Network    string `json:"network"`
			Strategies []struct {
				Params map[string]any `json:"params"`
			} `json:"strategies"`
		} `json:"space"`
	} `json:"data"`
}
type TokenIdAnswer struct {
	network       string
	tokensAddress []string
}

func GetTokenId(spaceID string) (TokenIdAnswer, error) {
	var result TokenIdAnswer
	var res SnapshotResponce

	// адрес для запроса
	apiURL := "https://hub.snapshot.org/graphql"

	// запроc
	query := map[string]string{
		"query": fmt.Sprintf(queryToken, spaceID),
	}
	// в json формат
	jsonData, err := json.Marshal(query)
	if err != nil {
		log.Println("JSON marshall error: ", err)
		return TokenIdAnswer{}, err
	}
	// запрос на snapshot
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Http error: ", err)
		return TokenIdAnswer{}, err
	}
	// особождение ресурсов
	defer resp.Body.Close()
	// вытаскиваем тело из ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read body error: ", err)
		return TokenIdAnswer{}, err
	}
	// парсим ответ в структуру
	if err := json.Unmarshal(body, &res); err != nil {
		log.Println("Unmarshal error:", err)
		return TokenIdAnswer{}, err
	}

	// кладем номер сети в структуру
	result.network = res.Data.Space.Network

	// кладем итоговые адреса токенов
	for _, s := range res.Data.Space.Strategies {
		// в params могут быть разные ключи: address, symbol, token и тд
		// проверяем address

		if addr, ok := s.Params["address"].(string); ok {
			result.tokensAddress = append(result.tokensAddress, addr)
		}
	}

	return result, nil
}
