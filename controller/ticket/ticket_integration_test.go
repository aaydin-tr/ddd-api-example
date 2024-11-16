package ticket

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	domain "github.com/aaydin-tr/gowit-case/domain/ticket"
	"github.com/aaydin-tr/gowit-case/domain/ticket/repository"
	"github.com/aaydin-tr/gowit-case/infrastructure/db/postgresql"
	"github.com/aaydin-tr/gowit-case/interface/http/request"
	"github.com/aaydin-tr/gowit-case/pkg/validator"
	service "github.com/aaydin-tr/gowit-case/service/ticket"
	"github.com/labstack/echo/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ticketTestSuite struct {
	suite.Suite
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	controller *TicketController
	sqlDB      *sql.DB
}

type ticketTestCase struct {
	name               string
	request            *string
	ticketID           int
	expectedStatusCode int
	expectedResponse   *string
}

var (
	createTicketTestCases = []ticketTestCase{
		{
			name:               "Create ticket successfully",
			request:            strToPointer(`{ "name": "example", "description": "sample description", "allocation": 100 }`),
			expectedResponse:   strToPointer(`{ "id": 1, "name": "example", "description": "sample description", "allocation": 100 }`),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Create ticket with invalid request (empty name)",
			request:            strToPointer(`{ "name": "", "description": "sample description", "allocation": 100 }`),
			expectedResponse:   strToPointer(`{ "message": "Validation error", "errors": [ { "failed_field": "Name", "tag": "required", "message": "This field is required" } ], "status": 400 }`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Create ticket with invalid request (empty desc)",
			request:            strToPointer(`{ "name": "example", "description": "", "allocation": 100 }`),
			expectedResponse:   strToPointer(`{ "message": "Validation error", "errors": [ { "failed_field": "Description", "tag": "required", "message": "This field is required" } ], "status": 400 }`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Create ticket with invalid request (empty allocation)",
			request:            strToPointer(`{ "name": "example", "description": "sample description", "allocation": 0 }`),
			expectedResponse:   strToPointer(`{ "message": "Validation error", "errors": [ { "failed_field": "Allocation", "tag": "required", "message": "This field is required" } ], "status": 400 }`),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	findByIDTicketTestCases = []ticketTestCase{
		{
			name:               "Get ticket successfully",
			request:            nil,
			ticketID:           1,
			expectedResponse:   strToPointer(`{ "id": 1, "name": "example", "description": "sample description", "allocation": 100 }`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Get non-existing ticket",
			request:            nil,
			ticketID:           2,
			expectedResponse:   strToPointer(`{ "message": "ticket not found", "status": 404, "errors": null }`),
			expectedStatusCode: http.StatusNotFound,
		},
	}

	purchasesTicketTestCases = []ticketTestCase{
		{
			name:               "Purchase ticket successfully",
			request:            strToPointer(`{ "quantity": 10, "user_id": "406c1d05-bbb2-4e94-b183-7d208c2692e1" }`),
			expectedResponse:   nil,
			expectedStatusCode: http.StatusOK,
			ticketID:           1,
		},
		{
			name:               "Purchase ticket with invalid request (empty quantity)",
			request:            strToPointer(`{ "quantity": 0, "user_id": "406c1d05-bbb2-4e94-b183-7d208c2692e1" }`),
			expectedResponse:   strToPointer(`{ "message": "Validation error", "errors": [ { "failed_field": "Quantity", "tag": "required", "message": "This field is required" } ], "status": 400 }`),
			expectedStatusCode: http.StatusBadRequest,
			ticketID:           1,
		},
		{
			name:               "Purchase ticket with invalid request (empty user_id)",
			request:            strToPointer(`{ "quantity": 10, "user_id": "" }`),
			expectedResponse:   strToPointer(`{ "message": "Validation error", "errors": [ { "failed_field": "UserID", "tag": "required", "message": "This field is required" } ], "status": 400 }`),
			expectedStatusCode: http.StatusBadRequest,
			ticketID:           1,
		},
		{
			name:               "Purchase ticket with invalid request (non-existing ticket)",
			request:            strToPointer(`{ "quantity": 10, "user_id": "406c1d05-bbb2-4e94-b183-7d208c2692e1" }`),
			expectedResponse:   strToPointer(`{ "message": "ticket not found", "status": 404, "errors": null }`),
			expectedStatusCode: http.StatusNotFound,
			ticketID:           2,
		},
		{
			name:               "Purchase ticket with invalid request (insufficient allocation)",
			request:            strToPointer(`{ "quantity": 1000, "user_id": "406c1d05-bbb2-4e94-b183-7d208c2692e1" }`),
			expectedResponse:   strToPointer(`{ "message": "insufficient allocation", "status": 422, "errors": null }`),
			expectedStatusCode: http.StatusUnprocessableEntity,
			ticketID:           1,
		},
	}
)

func (s *ticketTestSuite) SetupSuite() {
	var err error
	s.pool, err = dockertest.NewPool("")
	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}

	err = s.pool.Client.Ping()
	if err != nil {
		s.FailNow("Could not connect to Docker: %s", err.Error())
	}

	s.resource, err = s.pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "17",
			Env: []string{
				"POSTGRES_USER=postgres",
				"POSTGRES_PASSWORD=secret",
				"POSTGRES_DB=ticket",
			},
		},
		func(hostConfig *docker.HostConfig) {
			hostConfig.AutoRemove = true
			hostConfig.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)

	if err != nil {
		s.FailNow("Could not start resource: %s", err)
	}

	var dbClient *gorm.DB
	err = s.pool.Retry(func() error {
		var err error
		intPort, _ := strconv.Atoi(s.resource.GetPort("5432/tcp"))
		db, err := postgresql.NewPostgresDB("localhost", "postgres", "secret", "ticket", intPort)
		if err != nil {
			return err
		}

		dbClient = db
		return dbClient.AutoMigrate(&domain.Ticket{})
	})
	if err != nil {
		s.FailNow("Could not complete postgres migrations: %s", err)
	}

	sqlDB, err := dbClient.DB()
	if err != nil {
		s.Fail("Could not get sql.DB: %s", err)
	}

	s.sqlDB = sqlDB
	repos := repository.NewTicketRepository(dbClient)
	svc := service.NewTicketService(repos)
	controller := NewTicketController(svc)

	s.controller = controller
}

func (s *ticketTestSuite) TearDownSuite() {
	if err := s.sqlDB.Close(); err != nil {
		s.Fail("Could not close sql.DB: %s", err)
	}

	if err := s.pool.Purge(s.resource); err != nil {
		s.Fail("Could not purge resource: %s", err)
	}
}

func TestTicketTestSuite(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	suite.Run(t, &ticketTestSuite{})
}

func (s *ticketTestSuite) TestCreate() {
	for _, tc := range createTicketTestCases {
		s.T().Run(tc.name, func(t *testing.T) {
			api := echo.New()
			api.Validator = validator.New()
			api.POST("/tickets", s.controller.Create)

			req := httptest.NewRequest(http.MethodPost, "/tickets", strings.NewReader(pointerToStr(tc.request)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			api.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, *tc.expectedResponse, rec.Body.String())

			if rec.Code == http.StatusCreated {
				var ticket domain.Ticket
				var expectedTicket domain.TicketDTO
				err := json.Unmarshal([]byte(*tc.expectedResponse), &expectedTicket)
				if err != nil {
					t.Fatalf("could not unmarshal expected response: %s", err)
				}

				if err := s.sqlDB.QueryRow("SELECT id, name, description, allocation FROM tickets WHERE id = $1", expectedTicket.ID).Scan(&ticket.ID, &ticket.Name, &ticket.Description, &ticket.Allocation); err != nil {
					t.Fatalf("could not query ticket: %s", err)
				}

				assert.Equal(t, expectedTicket.ID, ticket.ID)
				assert.Equal(t, expectedTicket.Name, ticket.Name.GetValue())
				assert.Equal(t, expectedTicket.Description, ticket.Description.GetValue())
				assert.Equal(t, expectedTicket.Allocation, ticket.Allocation.GetValue())
			}

		})
	}
}

func (s *ticketTestSuite) TestFindByID() {
	for _, tc := range findByIDTicketTestCases {
		s.T().Run(tc.name, func(t *testing.T) {
			api := echo.New()
			api.Validator = validator.New()
			api.GET("/tickets/:id", s.controller.FindByID)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tickets/%d", tc.ticketID), nil)
			rec := httptest.NewRecorder()
			api.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, *tc.expectedResponse, rec.Body.String())
		})
	}
}

func (s *ticketTestSuite) TestPurchases() {
	for _, tc := range purchasesTicketTestCases {
		s.T().Run(tc.name, func(t *testing.T) {
			api := echo.New()
			api.Validator = validator.New()
			api.POST("/tickets/:id/purchases", s.controller.Purchases)

			var lastAllocationCount int
			s.sqlDB.QueryRow("SELECT allocation FROM tickets WHERE id = $1", tc.ticketID).Scan(&lastAllocationCount)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/tickets/%d/purchases", tc.ticketID), strings.NewReader(pointerToStr(tc.request)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			api.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedResponse != nil {
				assert.JSONEq(t, *tc.expectedResponse, rec.Body.String())
			} else {
				assert.Empty(t, rec.Body.String())
			}

			if rec.Code == http.StatusOK {
				var ticket domain.Ticket
				var req request.PurchaseTicketRequest
				err := json.Unmarshal([]byte(*tc.request), &req)
				if err != nil {
					t.Fatalf("could not unmarshal expected response: %s", err)
				}

				if err := s.sqlDB.QueryRow("SELECT id, name, description, allocation FROM tickets WHERE id = $1", tc.ticketID).Scan(&ticket.ID, &ticket.Name, &ticket.Description, &ticket.Allocation); err != nil {
					t.Fatalf("could not query ticket: %s", err)
				}

				assert.Equal(t, tc.ticketID, ticket.ID)
				assert.Equal(t, lastAllocationCount-req.Quantity, ticket.Allocation.GetValue())
			}

		})
	}
}

func strToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func pointerToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
