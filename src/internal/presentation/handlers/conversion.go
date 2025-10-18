package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
)

type ConversionHandlers struct {
	TelopStore   repositories.TelopStore
	AudioService repositories.AudioService
	wsService    *websocket.WebSocketHub
}

func (h *ConversionHandlers) PostConversionTelop(c *gin.Context) {
	var conv requests.ConversionRequest
	if err := c.ShouldBindJSON(&conv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	convEntity := conv.ToDomainConversion()

	h.TelopStore.SetConversionTelop(convEntity)
	telopMessage := h.TelopStore.GetCurrentTelopMessage()
	if telopMessage.IsSome() {
		h.wsService.PushTelop(telopMessage.Unwrap())
	}

	h.AudioService.UnMute()
	c.JSON(http.StatusOK, gin.H{"ok": true})

}
