package main

import (
	controller "github.com/aaydin-tr/gowit-case/controller/ticket"
	"github.com/aaydin-tr/gowit-case/domain/ticket"
	"github.com/aaydin-tr/gowit-case/infrastructure/db/postgresql"
	"github.com/aaydin-tr/gowit-case/interface/http"

	"github.com/aaydin-tr/gowit-case/domain/ticket/repository"
	"github.com/aaydin-tr/gowit-case/pkg/env"
	service "github.com/aaydin-tr/gowit-case/service/ticket"
)

func main() {
	config := env.ParseEnv()
	db, err := postgresql.NewPostgresDB(config.PostgresHost, config.PostgresUser, config.PostgresPassword, config.PostgresDB, config.PostgresPort)
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&ticket.Ticket{}); err != nil {
		panic(err)

	}

	repo := repository.NewTicketRepository(db)
	service := service.NewTicketService(repo)
	cont := controller.NewTicketController(service)

	http.NewEchoServer(cont, config.Port)
}
