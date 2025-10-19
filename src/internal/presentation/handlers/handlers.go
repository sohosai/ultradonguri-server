package handlers

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/file"
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

	// /force_mute
	r.POST("/force_mute", func(c *gin.Context) {
		var muteReq requests.MuteStateRequest //jsonを受け取るため
		if err := c.ShouldBindJSON(&muteReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newMuteState := muteReq.ToDomainMute() //domainの型に変換

		if err := h.AudioService.SetForceMute(newMuteState.IsMuted); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	//仕様からなくなっている
	// r.GET("/force_mute", func(c *gin.Context) {
	// 	state, err := h.AudioService.GetMute()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	newMuteState := responses.NewMuteStateResponse(state) //返すjsonに変換するための型変換

	// 	c.JSON(http.StatusOK, newMuteState)
	// })

	// /performances
	r.GET("/performances", func(c *gin.Context) {
		perfs, err := file.GetPerformances()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newPerfState := responses.NewPerformancesResponse(perfs)

		c.JSON(http.StatusOK, newPerfState)
	})

	// /performance
	r.POST("/performance/start", func(c *gin.Context) {
		var perf requests.PerformanceRequest
		if err := c.ShouldBindJSON(&perf); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//テロップは後で 受け取る型が変わっているため要変更
		// perfEntity := perf.ToDomainPerformance()

		// h.TelopStore.SetPerformanceTelop(perfEntity)
		// telopMessage := h.TelopStore.GetCurrentTelopMessage()
		// if telopMessage.IsSome() {
		// 	h.wsService.PushTelop(telopMessage.Unwrap())
		// }

		h.AudioService.SetIsConversion(false)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.POST("/performance/music", func(c *gin.Context) {
		var perf requests.MusicPostRequest
		if err := c.ShouldBindJSON(&perf); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		musicEntity := perf.ToDomainMusicPost()

		//テロップは後で　受け取る型が変わっているため要変更
		// h.TelopStore.SetPerformanceTelop(perfEntity)
		// telopMessage := h.TelopStore.GetCurrentTelopMessage()
		// if telopMessage.IsSome() {
		// 	h.wsService.PushTelop(telopMessage.Unwrap())
		// }

		//performance中しか/musicを呼べなくするなら、そのステートもいるかも
		//一旦簡易的にこちらでもisConersionをfalseにしておく
		h.AudioService.SetIsConversion(false)
		err := h.AudioService.SetShouldBeMuted(musicEntity.Music.ShouldBeMuted)
		if err != nil {
			//後でエラーを細かくする
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if musicEntity.Music.ShouldBeMuted {
			h.AudioService.SetMute(true)
		} else {
			h.AudioService.SetMute(false)
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// /conversion
	r.POST("/conversion/start", func(c *gin.Context) {
		var conv requests.ConversionRequest
		if err := c.ShouldBindJSON(&conv); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// convEntity := conv.ToDomainConversion()

		// h.TelopStore.SetConversionTelop(convEntity)
		// telopMessage := h.TelopStore.GetCurrentTelopMessage()
		// if telopMessage.IsSome() {
		// 	h.wsService.PushTelop(telopMessage.Unwrap())
		// }

		h.AudioService.SetIsConversion(true)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.POST("/conversion/cm-mode", func(c *gin.Context) {
		var conv requests.CMStateRequest
		if err := c.ShouldBindJSON(&conv); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		convEntity := conv.ToDomainCMState()

		//ここでのテロップはいらないかも
		// h.TelopStore.SetConversionTelop(convEntity)
		// telopMessage := h.TelopStore.GetCurrentTelopMessage()
		// if telopMessage.IsSome() {
		// 	h.wsService.PushTelop(telopMessage.Unwrap())
		// }

		if convEntity.IsCMMode {
			err := h.AudioService.SetCMScene()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		} else {
			err := h.AudioService.SetMute(false)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// /ws
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
