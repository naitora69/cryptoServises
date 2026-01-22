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

	"governance-indexer/internal/config"

	_ "github.com/lib/pq"
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
		log.Println(err)
	}

	// Подключения модулей
	repo := repository.NewRepository(db)
	index := indexer.NewIndexer(repo, cfg)
	tm := timer.NewTimer(index, cfg)

	// TODO: вместо это ерунды должен быть один поток очереди с ограничением количества запросов 60 в секунду
	//go tm.StartProposal()
	//go tm.StartSpace(false)
	go tm.StartVotes()

	// Создаем http сервер
	log.Println(fmt.Sprintf("Server started on: %s", cfg.Server.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), nil); err != nil {
		return
	}
}
