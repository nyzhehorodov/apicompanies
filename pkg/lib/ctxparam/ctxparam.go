// Package ctxparam provide helper functions to work with parameters incoming from the request url.
package ctxparam

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

// Package errors.
var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrConvertToString = errors.New("convert value to string fail")
)

// param is a key for use with context.WithValue.
// The param is required to isolate parameters
// from other data that can be in the context
// of the request.
type param string

func (p param) String() string { return "ctxparam:" + string(p) }

// WithValue returns a copy of parent context with the value.
func WithValue(parent context.Context, key string, val interface{}) context.Context {
	return context.WithValue(parent, param(key), val)
}

// Value returns an interface{} value from the context or ErrKeyNotFound.
func Value(ctx context.Context, key string) (interface{}, error) {
	val := ctx.Value(param(key))
	if val == nil {
		return nil, ErrKeyNotFound
	}
	return val, nil
}

// String returns a string value from the context or an error.
// If the key is not found returns ErrKeyNotFound.
// If the value is not a string returns ErrConvertToString.
func String(ctx context.Context, key string) (string, error) {
	val, err := Value(ctx, key)
	if err != nil {
		return "", err
	}
	s, ok := val.(string)
	if !ok {
		return "", ErrConvertToString
	}

	return s, nil
}

// Int returns an integer value from the context.
// It expect integer value written as string.
// If the key is not found returns ErrKeyNotFound.
// If the value is not an integer, it returns an error.
func Int(ctx context.Context, key string) (int, error) {
	val, err := String(ctx, key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("atoi: %s", err)
	}
	return i, nil
}
