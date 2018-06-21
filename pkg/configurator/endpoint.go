package configurator

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type getConfigRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type getConfigResponse struct {
	Data map[string]interface{}
	Err  error `json:"error,omitempty"`
}

func (r getConfigResponse) error() error { return r.Err }

func makeGetConfigEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getConfigRequest)
		data, err := s.GetConfig(req.Type, req.Data)
		return getConfigResponse{Data: data, Err: err}, nil
	}
}
