package ticket

import (
	"testing"

	"github.com/aaydin-tr/gowit-case/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewTicketDTOFromEntity(t *testing.T) {
	name, _ := valueobject.NewName("Test Ticket")
	description, _ := valueobject.NewDescription("This is a test ticket")
	allocation, _ := valueobject.NewAllocation(100)
	ticket := &Ticket{
		ID:          1,
		Name:        name,
		Description: description,
		Allocation:  allocation,
	}

	expected := &TicketDTO{
		ID:          1,
		Name:        name.GetValue(),
		Description: description.GetValue(),
		Allocation:  allocation.GetValue(),
	}

	result := NewTicketDTOFromEntity(ticket)
	assert.Equal(t, expected, result)
}
