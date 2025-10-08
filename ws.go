package main

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"example.com/donguri-back/spec"
	"github.com/gorilla/websocket"
)

// websocket の通信の集合
var WsConnections = struct {
	sync.Mutex
	conns map[*websocket.Conn]bool
}{conns: make(map[*websocket.Conn]bool)}

var telopQueue = make(chan spec.TelopMessage, 5)

func StartTelopWebsocketBroadcastWorker() {
	for telop := range telopQueue {
		WsConnections.Lock()

		conns := make([]*websocket.Conn, 0, len(WsConnections.conns))
		for conn, _ := range WsConnections.conns {
			conns = append(conns, conn)
		}

		WsConnections.Unlock()

		payload, err := json.Marshal(telop)

		if err != nil {
			slog.Error("Marsial error: " + err.Error())
		}

		for _, conn := range conns {
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			err := conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				slog.Error("Error sending message: " + err.Error())
				conn.Close()
				delete(WsConnections.conns, conn)
			}
		}
	}
}

func PushTelop(telop spec.TelopMessage) {
	telopQueue <- telop
}
