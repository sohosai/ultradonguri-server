package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/file"
)

// Handler は AudioService と TelopService を保持
type Handler struct {
	AudioService    service.AudioService
	TelopService    service.TelopService
	PerformanceRepo *file.PerformanceRepository
}

func NewHandler(audio service.AudioService, telop service.TelopService, repo *file.PerformanceRepository) *Handler {
	return &Handler{
		AudioService:    audio,
		TelopService:    telop,
		PerformanceRepo: repo,
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

		if err := h.AudioService.SetMute(muteReq.IsMuted); err != nil {
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
		perfs, err := h.PerformanceRepo.GetPerformances()
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
}
