package database

import "fmt"

const (
	DuplicateKeyPrefix = "duplicate key value violates unique constraint"
	NoRowsInResultSet  = "no rows in result set"
)

func NewInsertError(model string, err error) error {
	return fmt.Errorf("failed to insert table %s: %w", model, err)
}

func NewGetError(model string, err error) error {
	return fmt.Errorf("failed to query table %s: %w", model, err)
}
