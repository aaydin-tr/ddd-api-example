package service

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaydin-tr/ddd-api-example/domain/ticket"
	"github.com/aaydin-tr/ddd-api-example/interface/http/request"
	repository "github.com/aaydin-tr/ddd-api-example/mock/repository/ticket"
	"github.com/aaydin-tr/ddd-api-example/valueobject"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewTicketService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockTicketRepository(ctrl)
	service := NewTicketService(mockRepo)

	assert.NotNil(t, service)
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockTicketRepository(ctrl)
	service := NewTicketService(mockRepo)

	tests := []struct {
		name    string
		req     request.CreateTicketRequest
		mock    func()
		want    *ticket.TicketDTO
		wantErr bool
	}{
		{
			name: "success",
			req: request.CreateTicketRequest{
				Name:        "Test Ticket",
				Description: "Test Description",
				Allocation:  100,
			},
			mock: func() {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &ticket.TicketDTO{
				Name:        "Test Ticket",
				Description: "Test Description",
				Allocation:  100,
			},
			wantErr: false,
		},
		{
			name: "validation error",
			req: request.CreateTicketRequest{
				Name:        "",
				Description: "Test Description",
				Allocation:  10,
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repo error",
			req: request.CreateTicketRequest{
				Name:        "Test Ticket",
				Description: "Test Description",
				Allocation:  100,
			},
			mock: func() {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("repo error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.Create(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.Allocation, got.Allocation)
		})
	}
}
func TestService_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockTicketRepository(ctrl)
	service := NewTicketService(mockRepo)

	name, _ := valueobject.NewName("Test Ticket")
	description, _ := valueobject.NewDescription("Test Description")
	allocation, _ := valueobject.NewAllocation(100)
	tests := []struct {
		name    string
		id      int
		mock    func()
		want    *ticket.TicketDTO
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			mock: func() {
				mockRepo.EXPECT().FindByID(gomock.Any(), 1).Return(&ticket.Ticket{
					Name:        name,
					Description: description,
					Allocation:  allocation,
				}, nil)
			},
			want: &ticket.TicketDTO{
				Name:        name.GetValue(),
				Description: description.GetValue(),
				Allocation:  allocation.GetValue(),
			},
			wantErr: false,
		},
		{
			name: "not found error",
			id:   2,
			mock: func() {
				mockRepo.EXPECT().FindByID(gomock.Any(), 2).Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.FindByID(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.Allocation, got.Allocation)
		})
	}
}
func TestService_DecrementAllocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockTicketRepository(ctrl)
	service := NewTicketService(mockRepo)

	name, _ := valueobject.NewName("Test Ticket")
	description, _ := valueobject.NewDescription("Test Description")
	allocation, _ := valueobject.NewAllocation(100)

	tests := []struct {
		name     string
		ticketID int
		amount   int
		mock     func()
		wantErr  bool
	}{
		{
			name:     "success",
			ticketID: 1,
			amount:   50,
			mock: func() {
				mockDb, mock, _ := sqlmock.New()
				mock.ExpectBegin()
				mock.ExpectCommit()
				mock.ExpectRollback()
				dialector := postgres.New(postgres.Config{
					Conn:       mockDb,
					DriverName: "postgres",
				})
				db, _ := gorm.Open(dialector, &gorm.Config{})

				mockRepo.EXPECT().GetDB(gomock.Any()).Return(db)
				mockRepo.EXPECT().FindByIDForUpdate(gomock.Any(), 1, gomock.Any()).Return(&ticket.Ticket{
					Name:        name,
					Description: description,
					Allocation:  allocation,
				}, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "ticket not found error",
			ticketID: 2,
			amount:   50,
			mock: func() {
				mockDb, mock, _ := sqlmock.New()
				mock.ExpectBegin()
				mock.ExpectRollback()
				dialector := postgres.New(postgres.Config{
					Conn:       mockDb,
					DriverName: "postgres",
				})
				db, _ := gorm.Open(dialector, &gorm.Config{})

				mockRepo.EXPECT().GetDB(gomock.Any()).Return(db)
				mockRepo.EXPECT().FindByIDForUpdate(gomock.Any(), 2, gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:     "invalid amount error",
			ticketID: 1,
			amount:   150,
			mock: func() {
				mockDb, mock, _ := sqlmock.New()
				mock.ExpectBegin()
				mock.ExpectRollback()
				dialector := postgres.New(postgres.Config{
					Conn:       mockDb,
					DriverName: "postgres",
				})
				db, _ := gorm.Open(dialector, &gorm.Config{})

				mockRepo.EXPECT().GetDB(gomock.Any()).Return(db)
				mockRepo.EXPECT().FindByIDForUpdate(gomock.Any(), 1, gomock.Any()).Return(&ticket.Ticket{
					Name:        name,
					Description: description,
					Allocation:  allocation,
				}, nil)
			},
			wantErr: true,
		},
		{
			name:     "update error",
			ticketID: 1,
			amount:   50,
			mock: func() {
				mockDb, mock, _ := sqlmock.New()
				mock.ExpectBegin()
				mock.ExpectRollback()
				dialector := postgres.New(postgres.Config{
					Conn:       mockDb,
					DriverName: "postgres",
				})
				db, _ := gorm.Open(dialector, &gorm.Config{})

				mockRepo.EXPECT().GetDB(gomock.Any()).Return(db)
				mockRepo.EXPECT().FindByIDForUpdate(gomock.Any(), 1, gomock.Any()).Return(&ticket.Ticket{
					Name:        name,
					Description: description,
					Allocation:  allocation,
				}, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("update error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := service.DecrementAllocation(context.Background(), tt.ticketID, tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
