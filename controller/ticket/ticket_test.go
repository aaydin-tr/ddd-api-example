package ticket

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aaydin-tr/gowit-case/domain/ticket"
	mockservice "github.com/aaydin-tr/gowit-case/mock/service/ticket"

	"github.com/aaydin-tr/gowit-case/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTicketController_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockservice.NewMockTicketService(ctrl)
	controller := NewTicketController(mockService)

	e := echo.New()
	e.Validator = validator.New()

	tests := []struct {
		name         string
		requestBody  string
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			requestBody: `{
				"name": "Test Ticket",
				"description": "Test Description",
				"allocation": 100
			}`,
			mock: func() {
				mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&ticket.TicketDTO{
					Name:        "Test Ticket",
					Description: "Test Description",
					Allocation:  100,
				}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "bind error",
			requestBody: `{
				"name": "Test Ticket",
				"description": "Test Description",
				"allocation": "invalid"
			}`,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validation error",
			requestBody: `{
				"name": "",
				"description": "Test Description",
				"allocation": 100
			}`,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: `{
				"name": "Test Ticket",
				"description": "Test Description",
				"allocation": 100
			}`,
			mock: func() {
				mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tickets", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.mock()
			err := controller.Create(c)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
func TestTicketController_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockservice.NewMockTicketService(ctrl)
	controller := NewTicketController(mockService)

	e := echo.New()

	tests := []struct {
		name         string
		paramID      string
		mock         func()
		expectedCode int
	}{
		{
			name:    "success",
			paramID: "1",
			mock: func() {
				mockService.EXPECT().FindByID(gomock.Any(), 1).Return(&ticket.TicketDTO{
					ID:          1,
					Name:        "Test Ticket",
					Description: "Test Description",
					Allocation:  100,
				}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "id is required",
			paramID:      "",
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid id",
			paramID:      "invalid",
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "service error",
			paramID: "1",
			mock: func() {
				mockService.EXPECT().FindByID(gomock.Any(), 1).Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tickets/"+tt.paramID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			tt.mock()
			err := controller.FindByID(c)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
func TestTicketController_Purchases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockservice.NewMockTicketService(ctrl)
	controller := NewTicketController(mockService)

	e := echo.New()
	e.Validator = validator.New()

	tests := []struct {
		name         string
		paramID      string
		requestBody  string
		mock         func()
		expectedCode int
	}{
		{
			name:    "success",
			paramID: "1",
			requestBody: `{
				"quantity": 2,
				"user_id": "1250052d-c061-4a1f-81f0-d88af3dcb3d5"
			}`,
			mock: func() {
				mockService.EXPECT().DecrementAllocation(gomock.Any(), 1, 2).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "bind error",
			paramID:      "1",
			requestBody:  `{"quantity": "invalid", "user_id": "1250052d-c061-4a1f-81f0-d88af3dcb3d5"}`,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "validation error",
			paramID: "1",
			requestBody: `{
				"quantity": 0,
				"user_id": "invalid-uuid"
			}`,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "id is required",
			paramID:      "",
			requestBody:  `{"quantity": 2, "user_id": "1250052d-c061-4a1f-81f0-d88af3dcb3d5"}`,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid id",
			paramID:      "invalid",
			requestBody:  `{"quantity": 2, "user_id": "1250052d-c061-4a1f-81f0-d88af3dcb3d5"}`,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "service error",
			paramID: "1",
			requestBody: `{
				"quantity": 2,
				"user_id": "1250052d-c061-4a1f-81f0-d88af3dcb3d5"
			}`,
			mock: func() {
				mockService.EXPECT().DecrementAllocation(gomock.Any(), 1, 2).Return(errors.New("service error"))
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tickets/"+tt.paramID+"/purchases", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			tt.mock()
			err := controller.Purchases(c)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
