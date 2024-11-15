package response

import "github.com/labstack/echo/v4"

type ValidationMessage struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	Message string               `json:"message"`
	Errors  []*ValidationMessage `json:"errors"`
	Status  int                  `json:"status"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func NewErrorRespone(c echo.Context, err error, status int) error {
	if errResp, ok := err.(*ErrorResponse); ok {
		return c.JSON(errResp.Status, errResp)
	}

	return c.JSON(status, &ErrorResponse{
		Status:  status,
		Message: err.Error(),
	})
}
