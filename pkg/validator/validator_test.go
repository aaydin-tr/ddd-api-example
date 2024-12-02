package validator

import (
	"net/http"
	"testing"

	"github.com/aaydin-tr/ddd-api-example/interface/http/response"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `validate:"required"`
	UUID string `validate:"uuid4"`
}

func TestNew(t *testing.T) {
	v := New()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Validator)
}

func TestValidate(t *testing.T) {
	v := New()

	tests := []struct {
		input    interface{}
		expected error
	}{
		{
			input: &TestStruct{
				Name: "Test",
				UUID: "550e8400-e29b-41d4-a716-446655440000",
			},
			expected: nil,
		},
		{
			input: &TestStruct{
				Name: "",
				UUID: "invalid-uuid",
			},
			expected: &response.ErrorResponse{
				Message: "Validation error",
				Status:  http.StatusBadRequest,
				Errors: []*response.ValidationMessage{
					{
						FailedField: "Name",
						Tag:         "required",
						Message:     "This field is required",
					},
					{
						FailedField: "UUID",
						Tag:         "uuid4",
						Message:     "This field must be valid uuid4",
					},
				},
			},
		},
	}

	for _, test := range tests {
		err := v.Validate(test.input)
		if test.expected == nil {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
			assert.IsType(t, test.expected, err)
			assert.Equal(t, test.expected.(*response.ErrorResponse).Message, err.(*response.ErrorResponse).Message)
			assert.Equal(t, test.expected.(*response.ErrorResponse).Status, err.(*response.ErrorResponse).Status)
			assert.Equal(t, len(test.expected.(*response.ErrorResponse).Errors), len(err.(*response.ErrorResponse).Errors))
		}
	}
}
