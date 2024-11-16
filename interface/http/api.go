package http

import (
	"github.com/aaydin-tr/gowit-case/controller/ticket"
	"github.com/aaydin-tr/gowit-case/pkg/validator"

	_ "github.com/aaydin-tr/gowit-case/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewEchoServer(tickectController *ticket.TicketController, port string) {
	e := echo.New()
	e.Validator = validator.New()

	e.POST("/ticketsuser", tickectController.Create)
	e.GET("/tickets/:id", tickectController.FindByID)
	e.POST("/tickets/:id/purchases", tickectController.Purchases)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start("localhost:" + port))
}
