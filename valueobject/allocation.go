package valueobject

import (
	"database/sql/driver"
	"errors"
)

var (
	ErrInvalidAllocation = errors.New("invalid allocation")
)

type Allocation struct {
	value int
}

func NewAllocation(value int) (*Allocation, error) {
	if value < 0 {
		return nil, ErrInvalidAllocation
	}

	return &Allocation{value: value}, nil
}

func (a *Allocation) GetValue() int {
	return a.value
}

func (a *Allocation) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	price, ok := value.(*Allocation)
	if !ok {
		return false
	}

	return a.value == price.value

}

func (n *Allocation) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if _, ok := value.(int64); ok {
		n.value = int(value.(int64))
		return nil
	}

	if _, ok := value.(int32); ok {
		n.value = int(value.(int32))
		return nil
	}

	n.value = value.(int)
	return nil
}

func (n *Allocation) Value() (driver.Value, error) {
	return n.value, nil
}
