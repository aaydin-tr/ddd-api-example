package ticket

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aaydin-tr/gowit-case/interface/http/request"
	"github.com/aaydin-tr/gowit-case/interface/http/response"
	service "github.com/aaydin-tr/gowit-case/service/ticket"
	"github.com/labstack/echo/v4"
)

var (
	ErrIDIsRequired = errors.New("id is required")
)

type TicketController struct {
	service service.TicketService
}

func NewTicketController(service service.TicketService) *TicketController {
	return &TicketController{service: service}
}

func (t *TicketController) Create(c echo.Context) error {
	var req request.CreateTicketRequest
	if err := c.Bind(&req); err != nil {
		return response.NewErrorRespone(c, err, http.StatusBadRequest)
	}

	if err := c.Validate(req); err != nil {
		return response.NewErrorRespone(c, err, http.StatusBadRequest)
	}

	ticket, err := t.service.Create(c.Request().Context(), req)
	if err != nil {
		return response.NewErrorRespone(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, ticket)
}

func (t *TicketController) FindByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.NewErrorRespone(c, ErrIDIsRequired, http.StatusBadRequest)
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.NewErrorRespone(c, err, http.StatusBadRequest)
	}

	ticket, err := t.service.FindByID(c.Request().Context(), idInt)
	if err != nil {
		return response.NewErrorRespone(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, ticket)
}

func (t *TicketController) Purchases(c echo.Context) error {
	var req request.PurchaseTicketRequest
	if err := c.Bind(&req); err != nil {
		return response.NewErrorRespone(c, err, http.StatusBadRequest)
	}

	if err := c.Validate(req); err != nil {
		return response.NewErrorRespone(c, err, http.StatusBadRequest)
	}

	id := c.Param("id")
	if id == "" {
		return response.NewErrorRespone(c, ErrIDIsRequired, http.StatusBadRequest)
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.NewErrorRespone(c, err, http.StatusBadRequest)
	}

	err = t.service.DecrementAllocation(c.Request().Context(), idInt, req.Quantity)
	if err != nil {
		return response.NewErrorRespone(c, err, http.StatusUnprocessableEntity)
	}

	return c.NoContent(http.StatusOK)
}
