package service

import (
	"context"

	"github.com/aaydin-tr/gowit-case/domain/ticket"
	"github.com/aaydin-tr/gowit-case/domain/ticket/repository"
	"github.com/aaydin-tr/gowit-case/infrastructure/db"
	"github.com/aaydin-tr/gowit-case/interface/http/request"
)

type TicketService interface {
	Create(ctx context.Context, req request.CreateTicketRequest) (*ticket.TicketDTO, error)
	FindByID(ctx context.Context, id int) (*ticket.TicketDTO, error)
	DecrementAllocation(ctx context.Context, ticketID, amount int) error
}

type Service struct {
	repo repository.TicketRepository
}

func NewTicketService(repo repository.TicketRepository) TicketService {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req request.CreateTicketRequest) (*ticket.TicketDTO, error) {
	t, err := ticket.NewTicket(req.Name, req.Description, req.Allocation)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return ticket.NewTicketDTOFromEntity(t), nil
}

func (s *Service) FindByID(ctx context.Context, id int) (*ticket.TicketDTO, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return ticket.NewTicketDTOFromEntity(t), nil
}

func (s *Service) DecrementAllocation(ctx context.Context, ticketID, amount int) error {
	txManager := db.NewTransactionManager(s.repo.GetDB(ctx))
	tx, err := txManager.Begin(ctx)
	if err != nil {
		return err
	}

	t, err := s.repo.FindByIDForUpdate(ctx, ticketID, tx)
	if err != nil {
		txManager.Rollback(ctx)
		return err
	}

	err = t.DecrementAllocation(ctx, amount)
	if err != nil {
		txManager.Rollback(ctx)
		return err
	}

	err = s.repo.Update(ctx, t, tx)
	if err != nil {
		txManager.Rollback(ctx)
		return err
	}

	return txManager.Commit(ctx)
}
