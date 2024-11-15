package valueobject

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDescription(t *testing.T) {
	t.Run("should return error when description is empty", func(t *testing.T) {
		desc, err := NewDescription("")
		assert.Nil(t, desc)
		assert.Equal(t, ErrDescriptionCannotBeEmpty, err)
	})

	t.Run("should create description when value is provided", func(t *testing.T) {
		desc, err := NewDescription("Test description")
		assert.NotNil(t, desc)
		assert.Nil(t, err)
		assert.Equal(t, "Test description", desc.GetValue())
	})
}

func TestDescription_GetValue(t *testing.T) {
	desc := &Description{value: "Test description"}
	assert.Equal(t, "Test description", desc.GetValue())
}

func TestDescription_Equals(t *testing.T) {
	desc1, _ := NewDescription("Test description")
	desc2, _ := NewDescription("Test description")
	desc3, _ := NewDescription("Another description")

	assert.True(t, desc1.Equals(desc2))
	assert.False(t, desc1.Equals(desc3))
	assert.False(t, desc1.Equals(nil))
}

func TestDescription_Scan(t *testing.T) {
	desc := &Description{}
	err := desc.Scan("Scanned description")
	assert.Nil(t, err)
	assert.Equal(t, "Scanned description", desc.GetValue())
}

func TestDescription_Value(t *testing.T) {
	desc, _ := NewDescription("Test description")
	val, err := desc.Value()
	assert.Nil(t, err)
	assert.Equal(t, driver.Value("Test description"), val)
}
