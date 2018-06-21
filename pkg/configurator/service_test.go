package configurator

import (
	"reflect"
	"testing"

	"github.com/l-vitaly/imhio-task/pkg/catalog"
	"github.com/l-vitaly/imhio-task/pkg/value"
)

func TestGetConfigData(t *testing.T) {
	var (
		originValue = map[string]interface{}{
			"host": "localhost",
			"port": "12345",
		}
	)
	catalogs := &mockCatalogRespository{
		catalogs: map[string]*catalog.Catalog{
			"Develop.mr_robot": {
				ID:   1,
				Type: "Develop.mr_robot",
			},
		},
	}
	values := &mockValueRespository{
		values: map[int]map[string]*value.Value{
			1: {
				"Database.processing": {
					Data: originValue,
				},
			},
		},
	}

	s := NewService(catalogs, values)

	v, err := s.GetConfig("Develop.mr_robot", "Database.processing")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, originValue) {
		t.Errorf("v = %s; want = %s", v, originValue)
	}
}

type mockCatalogRespository struct {
	catalogs map[string]*catalog.Catalog
}

func (r *mockCatalogRespository) CreateSchemas() (int64, int64, error) {
	return 0, 0, nil
}

func (r *mockCatalogRespository) Find(t string) (*catalog.Catalog, error) {
	if c, ok := r.catalogs[t]; ok {
		return c, nil
	}
	return nil, catalog.ErrUnknown
}

type mockValueRespository struct {
	values map[int]map[string]*value.Value
}

func (r *mockValueRespository) CreateSchemas() (int64, int64, error) {
	return 0, 0, nil
}

func (r *mockValueRespository) Find(catalogID int, name string) (*value.Value, error) {
	if c, ok := r.values[catalogID]; ok {
		if v, ok := c[name]; ok {
			return v, nil
		}
	}
	return nil, value.ErrUnknown
}
