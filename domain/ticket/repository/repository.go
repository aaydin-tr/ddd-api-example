package repository

import (
	"context"
	"errors"

	"github.com/aaydin-tr/gowit-case/domain/ticket"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -destination=../../../mock/repository/ticket/ticket.go -package=repository github.com/aaydin-tr/gowit-case/domain/ticket/repository TicketRepository
type TicketRepository interface {
	GetDB(ctx context.Context) *gorm.DB
	Create(ctx context.Context, t *ticket.Ticket) error
	FindByID(ctx context.Context, id int) (*ticket.Ticket, error)
	FindByIDForUpdate(ctx context.Context, id int, tx *gorm.DB) (*ticket.Ticket, error)
	Update(ctx context.Context, ticket *ticket.Ticket, tx *gorm.DB) error
}

type Repository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &Repository{db: db}
}

func (r *Repository) GetDB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *Repository) Create(ctx context.Context, t *ticket.Ticket) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *Repository) FindByID(ctx context.Context, id int) (*ticket.Ticket, error) {
	var t ticket.Ticket
	err := r.db.WithContext(ctx).First(&t, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ticket.ErrTicketNotFound
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *Repository) FindByIDForUpdate(ctx context.Context, id int, tx *gorm.DB) (*ticket.Ticket, error) {
	var t ticket.Ticket
	err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&t, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ticket.ErrTicketNotFound
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *Repository) Update(ctx context.Context, t *ticket.Ticket, tx *gorm.DB) error {
	err := tx.Save(t).Error
	if err != nil {
		return err
	}

	return nil
}
