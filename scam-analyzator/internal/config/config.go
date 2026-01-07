package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetApiKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	apiKey := os.Getenv("ETHER_KEY")
	return apiKey
}
