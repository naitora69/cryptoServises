package main

import (
	"fmt"
	"governance-indexer/internal/indexer"
	"governance-indexer/internal/repository"
	"governance-indexer/internal/timer"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"

	"governance-indexer/internal/config"
)

func main() {
	// загружаем настройки из .env и yml файлов
	config.LoadEnv(".env")
	cfg := config.LoadConfig("configs/config.yml")

	// Конфигурация и подключение к PostgreSQL
	postgresConf := repository.PostgresConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     cfg.Database.Host,
		Port:     strconv.Itoa(cfg.Database.Port),
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}
	db, err := repository.NewPostgresDB(postgresConf)
	if err != nil {
		log.Fatal(err)
	}

	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("localhost:%s", cfg.Kafka.Port)),
		Topic:    "dao-indexer",
		Balancer: &kafka.LeastBytes{},
	}
	defer func(kafkaWriter *kafka.Writer) {
		err := kafkaWriter.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(kafkaWriter)

	// Подключения модулей
	repo := repository.NewRepository(db)
	index := indexer.NewIndexer(repo, cfg)
	tm := timer.NewTimer(index, cfg)
	go tm.StartProposal()

	// Создаем http сервер
	log.Println(fmt.Sprintf("Server started on: %s", cfg.Server.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), nil); err != nil {
		return
	}
}
