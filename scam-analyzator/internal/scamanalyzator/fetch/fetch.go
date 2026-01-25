package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/danila-kuryakin/cryptoServises/scam-analyzator/internal/config"
)

var queryToken = `
{ 
  space(id: "%s") {
    symbol
    network
    strategies {
      params 
	  network
    } 
  }
}`

type SnapshotResponce struct {
	Data struct {
		Space struct {
			Symbol     string `json:"symbol"`
			Network    string `json:"network"`
			Strategies []struct {
				Params  map[string]any `json:"params"`
				Network string         `json:"network"`
			} `json:"strategies"`
		} `json:"space"`
	} `json:"data"`
}
type MoralisResponce struct {
	Address          string `json:"address"`
	VerifiedContract bool   `json:"verified_contract"` // Проверен ли контракт в эксплорере
	PossibleSpam     bool   `json:"possible_spam"`     // Внутренний фильтр Moralis
}
type TokensInfo struct {
	Address string
	Network string
	Chain   string
}
type TokenIdAnswer struct {
	Symbol string
	Tokens []TokensInfo
}

func GetTokenId(spaceID string) (TokenIdAnswer, []error) {
	var allErrors []error
	var resSnapshot SnapshotResponce

	// получаем конфиг для наших api
	_, apiKey := config.GetApiKey()
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
		allErrors = append(allErrors, err)
		return TokenIdAnswer{}, allErrors
	}
	// запрос на snapshot
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Http error: ", err)
		allErrors = append(allErrors, err)
		return TokenIdAnswer{}, allErrors
	}
	// особождение ресурсов
	defer resp.Body.Close()
	// вытаскиваем тело из ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read body error: ", err)
		allErrors = append(allErrors, err)
		return TokenIdAnswer{}, allErrors
	}
	// парсим ответ в структуру
	if err := json.Unmarshal(body, &resSnapshot); err != nil {
		log.Println("Unmarshal error:", err)
		allErrors = append(allErrors, err)
		return TokenIdAnswer{}, allErrors
	}
	// кладем в срезы наши имена и номера сетей каждому api свое

	realSymbol := GetRealTokenSymbol(resSnapshot.Data.Space.Symbol) // получаем symbol пр-ва например у stgao.eth это STG
	tokensData := make([]TokensInfo, 0)

	for _, v := range resSnapshot.Data.Space.Strategies {
		chainName := GetChainName(v.Network)
		tokenFromSymbol, err := GetTokensFromSymbol(chainName, v.Network, realSymbol, apiKey)
		if err != nil {
			log.Println("Problem when GetTokensFromSymbol: ", err)
			allErrors = append(allErrors, err)
			continue
		}
		tokensData = append(tokensData, tokenFromSymbol...)
	}
	if len(allErrors) != 0 {
		return TokenIdAnswer{
			Symbol: realSymbol,
			Tokens: tokensData,
		}, allErrors
	}
	return TokenIdAnswer{
		Symbol: realSymbol,
		Tokens: tokensData,
	}, nil
}
func GetTokensFromSymbol(chainName string, networkID string, symbol string, apiKey string) ([]TokensInfo, error) {
	var resMoralis []MoralisResponce
	url := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/erc20/metadata/symbols?chain=%s&symbols=%s", chainName, symbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Problems with GetRequest: ", err)
		return nil, err
	}
	req.Header.Add("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Problems with HTTP client: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Problems when read body: ", err)
		return nil, err
	}

	if err := json.Unmarshal(body, &resMoralis); err != nil {
		log.Println("Moralis unmarshall error: ", err)
		return nil, err
	}
	result := make([]TokensInfo, 0, len(resMoralis))

	for _, v := range resMoralis {

		if v.VerifiedContract && !v.PossibleSpam {
			tmp := TokensInfo{
				Chain:   chainName,
				Address: v.Address,
				Network: networkID,
			}
			result = append(result, tmp)
		}
	}
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
func GetRealTokenSymbol(s string) string {
	// Регулярка для удаления префиксов ve, x, s, vl в начале строки
	re := regexp.MustCompile(`(?i)^(ve|s|vl|x|st|a|y)([A-Z0-9]{2,})`)
	matches := re.FindStringSubmatch(s)

	if len(matches) > 2 {
		return matches[2] // Возвращаем только основную часть, например STG
	}
	return s
}
