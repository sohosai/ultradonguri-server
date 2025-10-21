package handlers

import (
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
)

type Handler struct {
	SceneManager repositories.SceneManager
	TelopManager repositories.TelopManager
	wsService    *websocket.WebSocketHub
}

func NewHandler(scene repositories.SceneManager, telop repositories.TelopManager, wsHub *websocket.WebSocketHub) *Handler {
	return &Handler{
		SceneManager: scene,
		TelopManager: telop,
		wsService:    wsHub,
	}
}

func (h *Handler) Handle(r *gin.Engine) {

	healthHandler := HealthHandler{}

	muteHandler := MuteHandler{
		SceneManager: h.SceneManager,
		TelopManager: h.TelopManager,
	}

	performancesHandler := PerformancesHandler{}

	performanceHandler := PerformanceHandler{
		SceneManager: h.SceneManager,
		TelopManager: h.TelopManager,
		wsService:    h.wsService,
	}

	conversionHandlers := ConversionHandlers{
		SceneManager: h.SceneManager,
		TelopManager: h.TelopManager,
		wsService:    h.wsService,
	}

	copyrightHandler := CopyRightHandler{
		wsService: h.wsService,
	}

	websocketHandlers := WebsocketHandlers{
		SceneManager: h.SceneManager,
		TelopManager: h.TelopManager,
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

	r.POST("/display-copyright", copyrightHandler.PostDisplayCopyRight)

	r.GET("/ws", websocketHandlers.GetWebsocketConnection)
}
