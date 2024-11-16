package ticket

type TicketDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Allocation  int    `json:"allocation"`
} // @Name TicketDTO

func NewTicketDTOFromEntity(ticket *Ticket) *TicketDTO {
	return &TicketDTO{
		ID:          ticket.ID,
		Name:        ticket.Name.GetValue(),
		Description: ticket.Description.GetValue(),
		Allocation:  ticket.Allocation.GetValue(),
	}
}
