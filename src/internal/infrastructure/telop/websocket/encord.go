package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type TypedWebSocketResponse[T any] struct {
	Type string
	Data T
}

// telop送信用の型に変換
// 型が矛盾していないか
func (t TypedWebSocketResponse[T]) Encode() (WebSocketResponse, error) {
	rt := reflect.TypeOf(t.Data)
	log.Printf("%v", rt)
	expectedType, ok := typeRegistry[rt]
	if !ok {
		return WebSocketResponse{}, fmt.Errorf("unregistered type: %v", rt)
	}
	if t.Type != expectedType {
		return WebSocketResponse{}, fmt.Errorf("type mismatch: expected %s for %v, got %s", expectedType, rt, t.Type)
	}
	log.Printf("t.Data: %v", t.Data)
	b, err := json.Marshal(t.Data)
	log.Printf("b: %v", b)
	if err != nil {
		return WebSocketResponse{}, err
	}
	return WebSocketResponse{Type: t.Type, Data: b}, nil
}
