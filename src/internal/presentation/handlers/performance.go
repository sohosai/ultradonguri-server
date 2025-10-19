package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"

	"github.com/gin-gonic/gin"
)

type PerformanceHandler struct {
	AudioService repositories.AudioService
	TelopStore   repositories.TelopStore
	wsService    *websocket.WebSocketHub
}

// PostPerformanceStart godoc
// @Summary      start performance
// @Description  endpoint for start performance
// @Tags         performance
// @Accept       json
// @Produce      json
// @Param performanceStart body requests.PerformanceRequest true "post performance request"
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /performance/start [post]
func (h *PerformanceHandler) PostPerformanceStart(c *gin.Context) {
	var perf requests.PerformanceRequest
	if err := c.ShouldBindJSON(&perf); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
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
	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
}

// PostPerformanceMusic godoc
// @Summary      start performance music
// @Description  endpoint for start performance music
// @Tags         performance
// @Accept       json
// @Produce      json
// @Param performanceMusic body requests.MusicPostRequest true "post performance music"
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /performance/music [post]
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
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	if musicEntity.ShouldBeMuted {
		h.AudioService.SetMute(true)
	} else {
		h.AudioService.SetMute(false)
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
}
