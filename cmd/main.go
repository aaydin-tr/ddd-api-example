package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	controller "github.com/aaydin-tr/ddd-api-example/controller/ticket"
	"github.com/aaydin-tr/ddd-api-example/domain/ticket"
	"github.com/aaydin-tr/ddd-api-example/infrastructure/db/postgresql"
	"github.com/aaydin-tr/ddd-api-example/interface/http"

	"github.com/aaydin-tr/ddd-api-example/domain/ticket/repository"
	"github.com/aaydin-tr/ddd-api-example/pkg/env"
	service "github.com/aaydin-tr/ddd-api-example/service/ticket"
)

func main() {
	config := env.ParseEnv()
	db, sqlDB, err := postgresql.NewPostgresDB(config.PostgresHost, config.PostgresUser, config.PostgresPassword, config.PostgresDB, config.PostgresPort)
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&ticket.Ticket{}); err != nil {
		panic(err)
	}

	repo := repository.NewTicketRepository(db)
	service := service.NewTicketService(repo)
	cont := controller.NewTicketController(service)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	svc := http.NewEchoServer(cont, config.Host, config.Port)
	go svc.Start()

	<-ctx.Done()
	log.Println("Shutting down the server")
	if err := svc.Shutdown(); err != nil {
		panic(err)
	}

	log.Println("Shutting down the database")
	if err := sqlDB.Close(); err != nil {
		panic(err)
	}

	log.Println("Server and database are down")
	log.Println("Goodbye!")
}
