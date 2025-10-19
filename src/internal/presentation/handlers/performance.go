package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"

	"github.com/gin-gonic/gin"
)

type PerformanceHandler struct {
	AudioService repositories.AudioService
	TelopStore   repositories.TelopStore
	wsService    *websocket.WebSocketHub
}

func (h *PerformanceHandler) PostPerformanceStart(c *gin.Context) {
	var perf requests.PerformanceRequest
	if err := c.ShouldBindJSON(&perf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	perfEntity := perf.ToDomainPerformance()

	h.TelopStore.SetPerformanceTelop(perfEntity)
	telopMessage := h.TelopStore.GetCurrentTelopMessage()
	if telopMessage.IsSome() {
		resp, err := websocket.TypedWebSocketResponse[websocket.PerformanceStartData]{
			Type: websocket.TypePerformanceStart,
			Data: websocket.ToDataPerfStart(perfEntity), //ちゃんと、getの関数を書いて、telopClientから読むべきかも
		}.Encode()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.wsService.PushTelop(resp)
	}

	h.AudioService.SetIsConversion(false)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *PerformanceHandler) PostPerformanceMusic(c *gin.Context) {
	var music requests.MusicRequest
	if err := c.ShouldBindJSON(&music); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	musicEntity := music.ToDomainMusicPost()

	h.TelopStore.SetMusicTelop(musicEntity)
	telopMessage := h.TelopStore.GetCurrentTelopMessage()
	if telopMessage.IsSome() {
		resp, err := websocket.TypedWebSocketResponse[websocket.PerformanceMusicData]{
			Type: websocket.TypePerformanceMusic,
			Data: websocket.ToDataPerfMusic(musicEntity),
		}.Encode()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.wsService.PushTelop(resp)
	}

	//performance中しか/musicを呼べなくするなら、そのステートもいるかも
	//一旦簡易的にこちらでもisConersionをfalseにしておく
	h.AudioService.SetIsConversion(false)
	err := h.AudioService.SetShouldBeMuted(musicEntity.ShouldBeMuted)
	if err != nil {
		//後でエラーを細かくする
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if musicEntity.ShouldBeMuted {
		h.AudioService.SetMute(true)
	} else {
		h.AudioService.SetMute(false)
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
