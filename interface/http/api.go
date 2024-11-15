package http

import (
	"github.com/aaydin-tr/gowit-case/controller/ticket"
	"github.com/aaydin-tr/gowit-case/pkg/validator"

	"github.com/labstack/echo/v4"
)

func NewEchoServer(tickectController *ticket.TicketController, port string) {
	e := echo.New()
	e.Validator = validator.New()

	e.POST("/ticketsuser", tickectController.Create)
	e.GET("/tickets/:id", tickectController.FindByID)
	e.POST("/tickets/:id/purchases", tickectController.Purchases)
	// TODO host?
	e.Logger.Fatal(e.Start("localhost:" + port))
}
