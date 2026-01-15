package ageanalyze

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MoralisTransferResponse struct {
	Result []struct {
		BlockTimestamp string `json:"block_timestamp"`
	} `json:"result"`
}

func GetTokenAge(apiKey string, address string, chain string) (int64, error) {
	url := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/erc20/%s/transfers?chain=%s&order=ASC&limit=1", address, chain)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("moralis error: %s", string(body))
	}

	var data MoralisTransferResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	if len(data.Result) == 0 {
		return 0, fmt.Errorf("транзакции не найдены")
	}
	t, err := time.Parse(time.RFC3339, data.Result[0].BlockTimestamp)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}
