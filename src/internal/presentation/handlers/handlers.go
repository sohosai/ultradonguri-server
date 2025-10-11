package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/service"

	"github.com/gin-gonic/gin"
	. "github.com/sohosai/ultradonguri-server/internal/infrastructure/file"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop"
)

// Handler は AudioService と TelopService を保持
type Handler struct {
	AudioService service.AudioService
	TelopService service.TelopService
	wsService    *telop.WebSocketHub
	// PerformanceRepo *file.PerformanceRepository
}

func NewHandler(audio service.AudioService, telop service.TelopService, wsHub *telop.WebSocketHub) *Handler {
	return &Handler{
		AudioService: audio,
		TelopService: telop,
		wsService:    wsHub,
	}
}

func (h *Handler) Handle(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		message := "Hello World"
		c.IndentedJSON(http.StatusOK, message)
	})

	// /force_mute
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

		h.TelopService.SetPerformanceTelop(perf)

		h.wsService.PushTelop(entities.TelopMessage{
			Type:            entities.TelopTypePerformance,
			PerformanceData: &perf,
		})

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

		h.TelopService.SetConversionTelop(conv)

		h.AudioService.UnMute()
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	//テスト用WebSocketエンドポイント
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// CORS対策：どこからでも受け取る
			return true
		},
	}
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}
		// defer conn.Close()
		h.wsService.AddConnection(conn)

		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read error:", err)
				break
			}
			fmt.Println("Received:", string(msg))

			err = conn.WriteMessage(mt, []byte("Hello from server!"))
			if err != nil {
				fmt.Println("write error:", err)
				break
			}
		}
	})
}
