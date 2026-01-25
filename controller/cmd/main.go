package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/danila-kuryakin/cryptoServises/controller/internal/config"
	"github.com/danila-kuryakin/cryptoServises/controller/internal/controller"
	"github.com/danila-kuryakin/cryptoServises/controller/internal/repository"
	"github.com/danila-kuryakin/cryptoServises/controller/internal/service"

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
	log.Println(fmt.Sprintf("Kafka server: localhost:%s", cfg.Kafka.Port))

	// Подключения модулей
	repo := repository.NewRepository(db)
	services := service.NewService(repo, cfg)
	h := controller.NewController(services)
	go h.InitListener()
	go h.InitMessageController()

	// Создаем http сервер
	log.Println(fmt.Sprintf("Server started on: %s", cfg.Server.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), nil); err != nil {
		return
	}
}
