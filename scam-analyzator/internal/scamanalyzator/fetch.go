package scamanalyzator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scam-analyzator-service/internal/config"
)

var queryToken = `
{ 
  space(id: "%s") {
    symbol
    network
    strategies {
      params 
    } 
  }
}`

type SnapshotResponce struct {
	Data struct {
		Space struct {
			Symbol     string `json:"symbol"`
			Network    string `json:"network"`
			Strategies []struct {
				Params map[string]any `json:"params"`
			} `json:"strategies"`
		} `json:"space"`
	} `json:"data"`
}
type MoralisResponce struct {
	Address string `json:"address"`
}
type TokenIdAnswer struct {
	Network       string
	TokensAddress []string
}

func GetTokenId(spaceID string) (TokenIdAnswer, error) {
	var result TokenIdAnswer
	var resSnapshot SnapshotResponce
	var resMoralis []MoralisResponce

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
	if err := json.Unmarshal(body, &resSnapshot); err != nil {
		log.Println("Unmarshal error:", err)
		return TokenIdAnswer{}, err
	}

	_, apiKey := config.GetApiKey()
	// кладем номер сети в структуру
	result.Network = resSnapshot.Data.Space.Network
	symbol := resSnapshot.Data.Space.Symbol
	chain := GetChainName(result.Network)
	url := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/erc20/metadata/symbols?chain=%s&symbols=%s", chain, symbol)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-API-Key", apiKey)

	resp, _ = http.DefaultClient.Do(req)
	bodys, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(bodys, &resMoralis); err != nil {
		log.Println("Moralis unmarshall error: ", err)
	}
	bufferToStruct := make([]string, 0, len(resMoralis))

	for _, v := range resMoralis {
		bufferToStruct = append(bufferToStruct, v.Address)
	}
	result.TokensAddress = bufferToStruct
	return result, nil
}
func GetChainName(networkID string) string {
	switch networkID {
	case "1":
		return "eth"
	case "137":
		return "polygon"
	case "56":
		return "bsc"
	case "42161":
		return "arbitrum"
	case "10":
		return "optimism"
	case "250":
		return "fantom"
	case "43114":
		return "avalanche"
	default:
		return "eth"
	}
}
