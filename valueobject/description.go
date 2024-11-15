package valueobject

import (
	"database/sql/driver"
	"errors"
)

var (
	ErrDescriptionCannotBeEmpty = errors.New("Description cannot be empty")
)

type Description struct {
	value string
}

func NewDescription(value string) (*Description, error) {
	if value == "" {
		return nil, ErrDescriptionCannotBeEmpty
	}

	return &Description{value: value}, nil
}

func (n *Description) GetValue() string {
	return n.value
}

func (n *Description) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	description, ok := value.(*Description)
	if !ok {
		return false
	}

	return n.value == description.value
}

func (n *Description) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	n.value = value.(string)
	return nil
}

func (n *Description) Value() (driver.Value, error) {
	return n.value, nil
}
