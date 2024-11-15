package ticket_test

import (
	"context"
	"testing"

	"github.com/aaydin-tr/gowit-case/domain/ticket"
	"github.com/stretchr/testify/assert"
)

func TestDecrementAllocation(t *testing.T) {
	ctx := context.Background()
	t.Run("should decrement allocation successfully", func(t *testing.T) {
		firstAllocation := 10
		tk, err := ticket.NewTicket("Test Ticket", "Test Description", firstAllocation)
		assert.NoError(t, err)

		decrementAmount := 5
		err = tk.DecrementAllocation(ctx, decrementAmount)
		assert.NoError(t, err)
		assert.Equal(t, firstAllocation-decrementAmount, tk.Allocation.GetValue())
	})

	t.Run("should return error when new allocation is invalid", func(t *testing.T) {
		firstAllocation := 10
		tk, err := ticket.NewTicket("Test Ticket", "Test Description", firstAllocation)
		assert.NoError(t, err)

		decrementAmount := 15
		err = tk.DecrementAllocation(ctx, decrementAmount)
		assert.Error(t, err)
		assert.Equal(t, firstAllocation, tk.Allocation.GetValue())
	})

	t.Run("should return error when allocation hits zero", func(t *testing.T) {
		firstAllocation := 10
		tk, err := ticket.NewTicket("Test Ticket", "Test Description", firstAllocation)
		assert.NoError(t, err)

		decrementAmount := 10
		err = tk.DecrementAllocation(ctx, decrementAmount)
		assert.NoError(t, err)
		assert.Equal(t, 0, tk.Allocation.GetValue())

		err = tk.DecrementAllocation(ctx, 1)
		assert.Error(t, err)
	})
}
func TestNewTicket(t *testing.T) {
	t.Run("should create a new ticket successfully", func(t *testing.T) {
		name := "Test Ticket"
		description := "Test Description"
		allocation := 10

		tk, err := ticket.NewTicket(name, description, allocation)
		assert.NoError(t, err)
		assert.NotNil(t, tk)
		assert.Equal(t, name, tk.Name.GetValue())
		assert.Equal(t, description, tk.Description.GetValue())
		assert.Equal(t, allocation, tk.Allocation.GetValue())
	})

	t.Run("should return error when name is invalid", func(t *testing.T) {
		name := ""
		description := "Test Description"
		allocation := 10

		tk, err := ticket.NewTicket(name, description, allocation)
		assert.Error(t, err)
		assert.Nil(t, tk)
	})

	t.Run("should return error when description is invalid", func(t *testing.T) {
		name := "Test Ticket"
		description := ""
		allocation := 10

		tk, err := ticket.NewTicket(name, description, allocation)
		assert.Error(t, err)
		assert.Nil(t, tk)
	})

	t.Run("should return error when allocation is invalid", func(t *testing.T) {
		name := "Test Ticket"
		description := "Test Description"
		allocation := -1

		tk, err := ticket.NewTicket(name, description, allocation)
		assert.Error(t, err)
		assert.Nil(t, tk)
	})
}
