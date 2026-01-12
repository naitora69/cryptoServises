package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetApiKey() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	apiEhtKey := os.Getenv("ETHER_KEY")
	apiMorKey := os.Getenv("MORALILS_KEY")
	return apiEhtKey, apiMorKey
}
