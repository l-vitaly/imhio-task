package catalog

import "errors"

// ErrUnknown unknown catalog.
var ErrUnknown = errors.New("unknown catalog")

// Catalog is the configurator catalog.
type Catalog struct {
	ID   int
	Type string
}

// Repository provides access a catalog store.
type Repository interface {
	Find(t string) (*Catalog, error)
	CreateSchemas() (int64, int64, error)
}
