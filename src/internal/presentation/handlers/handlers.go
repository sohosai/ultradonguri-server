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

	conversionHandlers := ConversionHandlers{
		SceneManager: h.SceneManager,
	}

	r.GET("/health", healthHandler.GetHealth)

	conversionRoutes := r.Group("/conversion")
	{
		conversionRoutes.POST("/cm-mode", conversionHandlers.PostConversionCMMode)
	}

	sceneHandlers := SceneHandlers{
		SceneManager: h.SceneManager,
	}

	r.POST("/burari-scene", sceneHandlers.PostBurariScene)
	r.POST("/cm-scene", sceneHandlers.PostCmScene)
	r.POST("/normal-scene", sceneHandlers.PostNormalScene)
}
