package websocket

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type TypedWebSocketResponse[T any] struct {
	Type string
	Data T
}

func (t TypedWebSocketResponse[T]) Encode() (WebSocketResponse, error) {
	rt := reflect.TypeOf(t.Data)
	expectedType, ok := typeRegistry[rt]
	if !ok {
		return WebSocketResponse{}, fmt.Errorf("unregistered type: %v", rt)
	}
	if t.Type != expectedType {
		return WebSocketResponse{}, fmt.Errorf("type mismatch: expected %s for %v, got %s", expectedType, rt, t.Type)
	}
	b, err := json.Marshal(t.Data)
	if err != nil {
		return WebSocketResponse{}, err
	}
	return WebSocketResponse{Type: t.Type, Data: b}, nil
}
