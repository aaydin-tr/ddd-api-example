package ticket

import (
	"context"
	"errors"
	"time"

	"github.com/aaydin-tr/ddd-api-example/valueobject"
	"gorm.io/gorm"
)

var (
	ErrInsufficientAllocation = errors.New("insufficient allocation")
	ErrAllocationIsZero       = errors.New("allocation is zero")
	ErrNameIsRequired         = errors.New("name is required")
	ErrDescriptionIsRequired  = errors.New("description is required")
	ErrTicketNotFound         = errors.New("ticket not found")
)

type Ticket struct {
	ID          int                      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        *valueobject.Name        `json:"name" gorm:"not null;type:varchar(255)"`
	Description *valueobject.Description `json:"description" gorm:"not null;type:varchar(255)"`
	Allocation  *valueobject.Allocation  `json:"allocation" gorm:"not null;type:int;default:0"`
	CreatedAt   time.Time                `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt   time.Time                `json:"updated_at" gorm:"not null;default:current_timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt           `json:"deleted_at" gorm:"index"`
}

func (t *Ticket) TableName() string {
	return "tickets"
}

func (t *Ticket) DecrementAllocation(ctx context.Context, amount int) error {
	if t.Allocation.GetValue() == 0 || t.Allocation.GetValue() < amount {
		return ErrInsufficientAllocation
	}

	newAllocation, err := valueobject.NewAllocation(t.Allocation.GetValue() - amount)
	if err != nil {
		return err
	}

	t.Allocation = newAllocation
	return nil
}

func NewTicket(name string, description string, allocation int) (*Ticket, error) {
	ticketName, err := valueobject.NewName(name)
	if err != nil {
		return nil, err
	}

	ticketDescription, err := valueobject.NewDescription(description)
	if err != nil {
		return nil, err
	}

	ticketAllocation, err := valueobject.NewAllocation(allocation)
	if err != nil {
		return nil, err
	}

	return &Ticket{
		Name:        ticketName,
		Description: ticketDescription,
		Allocation:  ticketAllocation,
	}, nil
}
