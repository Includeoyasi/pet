package main

import (
	"context"
	"github.com/Includeoyasi/pet/api"
	"github.com/Includeoyasi/pet/internal/config"
	"log"
	"time"

	catpb "github.com/Includeoyasi/cat/pkg/cat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	catConn, err := connectCat(cfg.CatGrpcAddr)
	if err != nil {
		log.Fatalf("error connect cat service; err = %v", err)
	}
	defer func() { _ = catConn.Close() }()

	// Запуск Rest API сервера
	handler := api.NewHandler(catpb.NewCatClient(catConn))
	httpServer := new(api.Server)

	err = httpServer.Run(cfg.HttpPort, handler.InitRoutes())
	if err != nil {
		log.Fatalf("error start http server; err = %v", err)
		return
	}

	// Грейс-фул шатдаун
}

func connectCat(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
}
