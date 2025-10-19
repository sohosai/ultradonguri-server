package handlers

import (
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
)

type WebsocketHandlers struct {
	TelopManager repositories.TelopManager
	AudioService repositories.AudioService
	wsService    *websocket.WebSocketHub
}

func (h *WebsocketHandlers) GetWebsocketConnection(c *gin.Context) {
	wsConnection, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	slog.Info("New websocket connection is established")
	h.wsService.AddConnection(wsConnection)

	defer func() {
		h.wsService.RemoveConnection(wsConnection)
	}()

	for {
		mt, msg, err := wsConnection.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Println("Received:", string(msg))

		err = wsConnection.WriteMessage(mt, []byte("Hello from server!"))
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
