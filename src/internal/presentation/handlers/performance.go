package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"

	"github.com/gin-gonic/gin"
)

type PerformanceHandler struct {
	AudioService repositories.AudioService
	TelopStore   repositories.TelopStore
}

func (h *PerformanceHandler) PostPerformanceStart(c *gin.Context) {
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
}

func (h *PerformanceHandler) PostPerformanceMusic(c *gin.Context) {
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
}
