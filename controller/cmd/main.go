package main

import (
	"controller/internal/config"
	"controller/internal/controller"
	"controller/internal/repository"
	"controller/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
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
	log.Println(fmt.Sprintf("Kafka server: localhost:%s", cfg.Kafka.Port))

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("localhost:%s", cfg.Kafka.Port)},
		Topic:   "dao-indexer",
		GroupID: "dao-indexer",
		// читать с конца, если группа новая
		StartOffset: kafka.LastOffset,
	})
	defer func(kafkaReader *kafka.Reader) {
		err := kafkaReader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(kafkaReader)

	// Подключения модулей
	repo := repository.NewRepository(db)
	services := service.NewService(repo, cfg)
	h := controller.NewController(services)
	go h.InitListener()

	// Создаем http сервер
	log.Println(fmt.Sprintf("Server started on: %s", cfg.Server.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), nil); err != nil {
		return
	}
}
