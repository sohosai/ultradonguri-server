package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/file"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
)

type PerformanceHandler struct {
	AudioService repositories.AudioService
	TelopStore   repositories.TelopStore
	wsService    *websocket.WebSocketHub
}

func (h *PerformanceHandler) GetPerformances(c *gin.Context) {
	perfs, err := file.GetPerformances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newPerfState := responses.NewPerformancesResponse(perfs)
	c.JSON(http.StatusOK, newPerfState)
}

func (h *PerformanceHandler) PostPerformance(c *gin.Context) {
	var perf requests.PerformancePostRequest
	if err := c.ShouldBindJSON(&perf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	perfEntity := perf.ToDomainPerformancePost()

	h.TelopStore.SetPerformanceTelop(perfEntity)
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
}
