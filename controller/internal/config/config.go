package config

import (
	"bufio"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config — корневая конфигурация приложения
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Proposal Proposal       `yaml:"proposal"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}
type DatabaseConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Name    string `yaml:"dbname"`
	SSLMode string `yaml:"sslmode"`
}

type Proposal struct {
	NumberRecords int `yaml:"number_records"` // Количество записей в ответе
	TimeRequest   int `yaml:"time_request"`   // Время между запросами
}

type KafkaConfig struct {
	Port string `yaml:"port"`
}

// LoadConfig загружает конфигурацию из YAML файла
func LoadConfig(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}

// LoadEnv загружает .env файл в окружение. Строки .env файла:
// DB_USER - Имя пользователя
// DB_PASSWORD - Пароль
func LoadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// пропускаем комментарии и пустые строки
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value)
	}
}
