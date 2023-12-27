package main

import (
	"github.com/Includeoyasi/pet/api"
	"github.com/Includeoyasi/pet/internal/config"
	"log"
)

func main() {

	// Загрузка конфигурации
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("error init config; err = %v", err)
	}

	// Инициализация соединения с Redis

	// Инициализация соединения с Postgres

	// Инициализация брокера Kafka

	// Запуск Rest API сервера
	handler := new(api.Handler)
	httpServer := new(api.Server)

	err = httpServer.Run(cfg.HttpPort, handler.InitRoutes())
	if err != nil {
		log.Fatalf("error start http server; err = %v", err)
		return
	}

	// Грейс-фул шатдаун
}
