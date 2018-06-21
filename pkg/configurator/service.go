package configurator

import (
	"github.com/l-vitaly/imhio-task/pkg/catalog"
	"github.com/l-vitaly/imhio-task/pkg/value"
)

// Service is the interface that provides configuration methods.
type Service interface {
	// GetConfig returns a configuration data.
	GetConfig(t, name string) (map[string]interface{}, error)
}

type service struct {
	catalogs catalog.Repository
	values   value.Repository
}

func (s *service) GetConfig(t, name string) (map[string]interface{}, error) {
	c, err := s.catalogs.Find(t)
	if err != nil {
		return nil, err
	}
	v, err := s.values.Find(c.ID, name)
	if err != nil {
		return nil, err
	}
	return v.Data, nil
}

// NewService creates a configuration service.
func NewService(
	catalogs catalog.Repository,
	values value.Repository,
) Service {
	return &service{
		catalogs: catalogs,
		values:   values,
	}
}
