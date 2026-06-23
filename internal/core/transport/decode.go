package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Validator interface {
	Valid(ctx context.Context) map[string]string
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		if err == io.EOF {
			return v, fmt.Errorf("request body cannot be empty")
		}
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func DecodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	v, err := Decode[T](r)
	if err != nil {
		return v, nil, err
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %w", v, ErrValidation)
	}
	return v, nil, nil
}
