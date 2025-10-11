package handlers

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"

	"github.com/gin-gonic/gin"
	. "github.com/sohosai/ultradonguri-server/internal/infrastructure/file"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
)

type Handler struct {
	AudioService repositories.AudioService
	TelopStore   repositories.TelopStore
	wsService    *websocket.WebSocketHub
}

func NewHandler(audio repositories.AudioService, telop repositories.TelopStore, wsHub *websocket.WebSocketHub) *Handler {
	return &Handler{
		AudioService: audio,
		TelopStore:   telop,
		wsService:    wsHub,
	}
}

func (h *Handler) Handle(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		message := "Hello World"
		c.IndentedJSON(http.StatusOK, message)
	})

	r.POST("/force_mute", func(c *gin.Context) {
		var muteReq entities.MuteState
		if err := c.ShouldBindJSON(&muteReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := h.AudioService.SetMute(muteReq.Is_Muted); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.GET("/force_mute", func(c *gin.Context) {
		state, err := h.AudioService.GetMute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, state)
	})

	// /performances
	r.GET("/performances", func(c *gin.Context) {
		perfs, err := GetPerformances()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, perfs)
	})

	// /performance
	r.POST("/performance", func(c *gin.Context) {
		var perf entities.PerformancePost
		if err := c.ShouldBindJSON(&perf); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.TelopStore.SetPerformanceTelop(perf)
		telopMessage := h.TelopStore.GetCurrentTelopMessage()
		if telopMessage.IsSome() {
			h.wsService.PushTelop(telopMessage.Unwrap())
		}

		if perf.Music.ShouldBeMuted {
			h.AudioService.Mute()
		} else {
			h.AudioService.UnMute()
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// /conversion
	r.POST("/conversion", func(c *gin.Context) {
		var conv entities.ConversionPost
		if err := c.ShouldBindJSON(&conv); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.TelopStore.SetConversionTelop(conv)
		telopMessage := h.TelopStore.GetCurrentTelopMessage()
		if telopMessage.IsSome() {
			h.wsService.PushTelop(telopMessage.Unwrap())
		}

		h.AudioService.UnMute()
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.GET("/ws", func(c *gin.Context) {
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
	})
}
