package request

type PurchaseTicketRequest struct {
	Quantity int    `json:"quantity" validate:"required,gte=1"`
	UserID   string `json:"user_id" validate:"required,uuid4"`
} // @Name PurchaseTicketRequest

type CreateTicketRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Allocation  int    `json:"allocation" validate:"required,gte=1"`
} // @Name CreateTicketRequest
