package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"petProject/internal/sender"

	"petProject/internal/handler"
	"petProject/internal/repository"
	"petProject/internal/service"
	"petProject/internal/syncer"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := repository.Config{
		Login:      "krasnovkd21",
		Password:   "yes",
		Database:   "petProjectDB",
		MasterHost: "db",
	}

	client, err := repository.NewClient(ctx, cfg.MasterHost, cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer client.Close()

	db := client.DB()
	repo := repository.NewTaskRepository(db)
	svc := service.NewTaskService(repo)
	taskHandler := handler.NewTaskHandler(svc)

	r := chi.NewRouter()
	taskHandler.RegisterRoutes(r)

	sender := sender.NewHTTPSender()
	go syncer.NewSyncer(svc, sender).Run(ctx)

	//MARK: - Запуск HTTP-сервера
	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
