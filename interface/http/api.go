package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/aaydin-tr/gowit-case/controller/ticket"
	"github.com/aaydin-tr/gowit-case/pkg/validator"

	_ "github.com/aaydin-tr/gowit-case/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type EchoServer struct {
	controller *ticket.TicketController
	host       string
	port       string

	e *echo.Echo
}

func NewEchoServer(tickectController *ticket.TicketController, host, port string) *EchoServer {
	svc := &EchoServer{
		controller: tickectController,
		host:       host,
		port:       port,
	}

	e := echo.New()
	e.Validator = validator.New()

	svc.e = e

	return svc
}

func (s *EchoServer) Start() {
	s.e.POST("/ticketsuser", s.controller.Create)
	s.e.GET("/tickets/:id", s.controller.FindByID)
	s.e.POST("/tickets/:id/purchases", s.controller.Purchases)
	s.e.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := s.e.Start(s.host + ":" + s.port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.e.Logger.Fatal(err)
	}
}

func (s *EchoServer) Shutdown() error {
	return s.e.Shutdown(context.Background())
}
