package valueobject

import (
	"database/sql/driver"
	"errors"
)

var (
	ErrNameCannotBeEmpty = errors.New("Name cannot be empty")
)

type Name struct {
	value string
}

func NewName(value string) (*Name, error) {
	if value == "" {
		return nil, ErrNameCannotBeEmpty
	}

	return &Name{value: value}, nil
}

func (n *Name) GetValue() string {
	return n.value
}

func (n *Name) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	name, ok := value.(*Name)
	if !ok {
		return false
	}

	return n.value == name.value
}

func (n *Name) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	n.value = value.(string)
	return nil
}

func (n *Name) Value() (driver.Value, error) {
	return n.value, nil
}
