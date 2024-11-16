package ticket

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aaydin-tr/gowit-case/domain/ticket"
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

// Create godoc
// @Summary      Create a new ticket
// @Description  Create a new ticket
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        ticket body request.CreateTicketRequest true "ticket"
// @Success      200  {object}  ticket.TicketDTO
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /ticketsuser [post]
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
		return response.NewErrorRespone(c, err, http.StatusUnprocessableEntity)
	}

	return c.JSON(http.StatusCreated, ticket)
}

// FindByID godoc
// @Summary      Find ticket by ID
// @Description  Find ticket by ID
// @Tags         tickets
// @Produce      json
// @Param        id path int true "ticket ID"
// @Success      200  {object}  ticket.TicketDTO
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /tickets/{id} [get]
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
		return response.NewErrorRespone(c, err, http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, ticket)
}

// Purchases godoc
// @Summary      Purchase tickets
// @Description  Purchase tickets
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        id path int true "ticket ID"
// @Param        purchase body request.PurchaseTicketRequest true "purchase"
// @Success      200  {object}  response.EmptyBody "No content"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Router       /tickets/{id}/purchases [post]
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
	if errors.Is(err, ticket.ErrTicketNotFound) {
		return response.NewErrorRespone(c, err, http.StatusNotFound)
	}

	if err != nil {
		return response.NewErrorRespone(c, err, http.StatusUnprocessableEntity)
	}

	return c.NoContent(http.StatusOK)
}
