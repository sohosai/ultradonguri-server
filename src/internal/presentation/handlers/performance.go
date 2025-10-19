package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"

	"github.com/gin-gonic/gin"
)

type PerformanceHandler struct {
	AudioService repositories.AudioService
	TelopStore   repositories.TelopStore
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

	//テロップは後で 受け取る型が変わっているため要変更
	// perfEntity := perf.ToDomainPerformance()

	// h.TelopStore.SetPerformanceTelop(perfEntity)
	// telopMessage := h.TelopStore.GetCurrentTelopMessage()
	// if telopMessage.IsSome() {
	// 	h.wsService.PushTelop(telopMessage.Unwrap())
	// }

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
	var perf requests.MusicPostRequest
	if err := c.ShouldBindJSON(&perf); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
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
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	if musicEntity.Music.ShouldBeMuted {
		h.AudioService.SetMute(true)
	} else {
		h.AudioService.SetMute(false)
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
}
