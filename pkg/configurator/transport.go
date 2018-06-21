package configurator

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/l-vitaly/imhio-task/pkg/catalog"
	"github.com/l-vitaly/imhio-task/pkg/value"
)

// MakeHandler returns a handler for the configurator service.
func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getConfigHandler := kithttp.NewServer(
		makeGetConfigEndpoint(s),
		decodeGetConfigRequest,
		encodeGetConfigResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/cfg", getConfigHandler).Methods("POST")

	return r
}

var errBadRoute = errors.New("bad route")

func decodeGetConfigRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeGetConfigResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response.(getConfigResponse).Data)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case value.ErrUnknown, catalog.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
