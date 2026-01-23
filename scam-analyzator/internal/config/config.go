package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetApiKey() (string, string) {
	// Пытаемся загрузить .env только если файл существует (для локальной разработки)
	// В Docker используем только переменные окружения
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("Ошибка загрузки .env файла")
		}
	}

	apiEhtKey := os.Getenv("ETHER_KEY")
	apiMorKey := os.Getenv("MORALILS_KEY")

	if apiEhtKey == "" {
		log.Println("Предупреждение: ETHER_KEY не установлен")
	}
	if apiMorKey == "" {
		log.Println("Предупреждение: MORALIS_KEY не установлен")
	}

	return apiEhtKey, apiMorKey
}
