package handlers

import (
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	SceneManager repositories.SceneManager
}

func NewHandler(scene repositories.SceneManager) *Handler {
	return &Handler{
		SceneManager: scene,
	}
}

func (h *Handler) Handle(r *gin.Engine) {

	healthHandler := HealthHandler{}

	muteHandler := MuteHandler{
		SceneManager: h.SceneManager,
	}

	performancesHandler := PerformancesHandler{}

	conversionHandlers := ConversionHandlers{
		SceneManager: h.SceneManager,
	}

	r.GET("/health", healthHandler.GetHealth)

	r.POST("/force_mute", muteHandler.PostForceMuted)
	r.GET("/performances", performancesHandler.GetPerformances)

	conversionRoutes := r.Group("/conversion")
	{
		conversionRoutes.POST("/cm-mode", conversionHandlers.PostConversionCMMode)
	}

}
