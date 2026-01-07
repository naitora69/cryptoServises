package scamanalyzator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type EtherscanResponse struct {
	Status  string          `json:"status"` // "1" = OK
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}
type Transaction struct {
	TimeStamp string `json:"timeStamp"`
}

func GetTokenAge(apiKey string, address string, chainID string) (int64, error) {
	// адрес для запроса
	apiUrl := fmt.Sprintf(
		"https://api.etherscan.io/v2/api?chainid=%s&module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=1&sort=asc&apikey=%s",
		chainID, address, apiKey)
	// chainId - номер сети, adress - адрес токена, apiKey - ключ с моего аккаунта в etherscan

	// отправляем запрос на сервер
	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Println("HTTP GET erorr: ", err)
		return 0, err
	}

	defer resp.Body.Close()

	var res EtherscanResponse

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Decode error: ", err)
		return 0, err
	}
	// статус не еденица в ответе - ошибка
	if res.Status != "1" {
		log.Println("Status err", fmt.Errorf("etherscan v2 error: %s (result: %s)", res.Message, string(res.Result)))
		return 0, fmt.Errorf("etherscan v2 error: %s (result: %s)", res.Message, string(res.Result))
	}

	var txs []Transaction
	if err := json.Unmarshal(res.Result, &txs); err != nil {
		log.Println("Transaction parsing: ", fmt.Errorf("error parsing transactions: %w", err))
		return 0, fmt.Errorf("error parsing transactions: %w", err)
	}

	if len(txs) == 0 {
		log.Println("Transaction error: ", fmt.Errorf("транзакции не найдены"))
		return 0, fmt.Errorf("транзакции не найдены")
	}
	// ответ в Unix Timestamp(1 января 1970)
	var ts int64
	fmt.Sscanf(txs[0].TimeStamp, "%d", &ts)

	return ts, nil
}
