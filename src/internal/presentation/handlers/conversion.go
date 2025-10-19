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

func (h *ConversionHandlers) PostConversionStart(c *gin.Context) {
	var conv requests.ConversionRequest
	if err := c.ShouldBindJSON(&conv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	convEntity := conv.ToDomainConversion()

	h.TelopStore.SetConversionTelop(convEntity)
	telopMessage := h.TelopStore.GetCurrentTelopMessage()
	if telopMessage.IsSome() {
		resp, err := websocket.TypedWebSocketResponse[websocket.ConversionStartData]{
			Type: websocket.TypeConversionStart,
			Data: websocket.ToDataConvStart(convEntity),
		}.Encode()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.wsService.PushTelop(resp)
	}

	h.AudioService.SetIsConversion(true)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ConversionHandlers) PostConversionCMMode(c *gin.Context) {
	var conv requests.CMStateRequest
	if err := c.ShouldBindJSON(&conv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	convEntity := conv.ToDomainCMState()

	resp, err := websocket.TypedWebSocketResponse[websocket.ConversionCmModeData]{
		Type: websocket.TypeConversionCmMode,
		Data: websocket.ToDataConvCmMode(convEntity),
	}.Encode()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.wsService.PushTelop(resp)

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
}
