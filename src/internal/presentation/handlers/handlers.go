package handlers

import (
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"

	"github.com/gin-gonic/gin"
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

	healthHandler := HealthHandler{}

	muteHandler := MuteHandler{
		AudioService: h.AudioService,
	}

	performancesHandler := PerformancesHandler{}

	performanceHandler := PerformanceHandler{
		AudioService: h.AudioService,
		TelopStore:   h.TelopStore,
		wsService:    h.wsService,
	}

	conversionHandlers := ConversionHandlers{
		AudioService: h.AudioService,
		TelopStore:   h.TelopStore,
		wsService:    h.wsService,
	}

	websocketHandlers := WebsocketHandlers{
		AudioService: h.AudioService,
		TelopStore:   h.TelopStore,
		wsService:    h.wsService,
	}

	r.GET("/health", healthHandler.GetHealth)

	r.POST("/force_mute", muteHandler.PostForceMuted)
	r.GET("/performances", performancesHandler.GetPerformances)

	performanceRoutes := r.Group("/performance")
	{
		performanceRoutes.POST("/start", performanceHandler.PostPerformanceStart)
		performanceRoutes.POST("/music", performanceHandler.PostPerformanceMusic)
	}

	conversionRoutes := r.Group("/conversion")
	{
		conversionRoutes.POST("/start", conversionHandlers.PostConversionStart)
		conversionRoutes.POST("/cm-mode", conversionHandlers.PostConversionCMMode)
	}

	r.GET("/ws", websocketHandlers.GetWebsocketConnection)
}
