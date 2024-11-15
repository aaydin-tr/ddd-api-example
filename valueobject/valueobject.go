package valueobject

import "database/sql/driver"

type ValueObject interface {
	Equals(value ValueObject) bool
	Scan(value interface{}) error
	Value() (driver.Value, error)
}
