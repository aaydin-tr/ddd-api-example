package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewName(t *testing.T) {
	t.Run("should return error when name is empty", func(t *testing.T) {
		desc, err := NewName("")
		assert.Nil(t, desc)
		assert.Equal(t, ErrNameCannotBeEmpty, err)
	})

	t.Run("should create name when value is provided", func(t *testing.T) {
		desc, err := NewName("Name")
		assert.NotNil(t, desc)
		assert.Nil(t, err)
		assert.Equal(t, "Name", desc.GetValue())
	})
}

func TestName_GetValue(t *testing.T) {
	name := &Name{value: "John Doe"}
	assert.Equal(t, "John Doe", name.GetValue())
}

func TestName_Equals(t *testing.T) {
	name1 := &Name{value: "John Doe"}
	name2 := &Name{value: "John Doe"}
	name3 := &Name{value: "Jane Doe"}

	assert.True(t, name1.Equals(name2))
	assert.False(t, name1.Equals(name3))
	assert.False(t, name1.Equals(nil))
}

func TestName_ScanAndValue(t *testing.T) {
	name := &Name{}

	err := name.Scan("John Doe")
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", name.value)

	val, err := name.Value()
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", val)

	err = name.Scan(nil)
	assert.NoError(t, err)
}
