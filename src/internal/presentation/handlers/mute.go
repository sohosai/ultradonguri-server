package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"

	"github.com/gin-gonic/gin"
)

type MuteHandler struct {
	AudioService repositories.AudioService
}

func (h *MuteHandler) PostForceMuted(c *gin.Context) {
	var muteReq requests.MuteStateRequest //jsonを受け取るため
	if err := c.ShouldBindJSON(&muteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newMuteState := muteReq.ToDomainMute() //domainの型に変換

	if err := h.AudioService.SetMute(newMuteState.IsMuted); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
