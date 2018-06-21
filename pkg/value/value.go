package value

import "errors"

// ErrUnknown unknown value.
var ErrUnknown = errors.New("unknown value")

// Value is the configuration value.
type Value struct {
	ID        int
	CatalogID int
	Name      string
	Data      map[string]interface{}
}

// Repository provides access a value store.
type Repository interface {
	Find(catalogID int, name string) (*Value, error)
	CreateSchemas() (int64, int64, error)
}
