package util

import (
	"bytes"
	"encoding/json"
)

type Option[T any] struct {
	Valid bool
	V     T
}

func Some[T any](v T) Option[T] { return Option[T]{Valid: true, V: v} }
func None[T any]() Option[T]    { return Option[T]{} }

func (o Option[T]) IsSome() bool { return o.Valid }
func (o Option[T]) IsNone() bool { return !o.Valid }
func (o Option[T]) UnwrapOr(def T) T {
	if o.Valid {
		return o.V
	}
	return def
}
func (o Option[T]) Unwrap() T {
	if o.Valid {
		return o.V
	}

	panic("Failed to unwrap.")
}

// JSON: None => null, Some(x) => x
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(o.V)
}
func (o *Option[T]) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		o.Valid = false
		var zero T
		o.V = zero
		return nil
	}
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	o.Valid, o.V = true, v
	return nil
}
