package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"
)

type ConversionHandlers struct {
	TelopStore   repositories.TelopStore
	AudioService repositories.AudioService
	wsService    *websocket.WebSocketHub
}

// PostConversionStart godoc
// @Summary      conversion start
// @Description  endpoint for start conversion
// @Tags         conversion
// @Accept       json
// @Produce      json
// @Param conversionStart body requests.ConversionRequest true "post conversion request"
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /conversion/start [post]
func (h *ConversionHandlers) PostConversionStart(c *gin.Context) {
	var conv requests.ConversionRequest
	if err := c.ShouldBindJSON(&conv); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	convEntity := conv.ToDomainConversion()

	h.TelopStore.SetConversionTelop(convEntity)
	telopMessage := h.TelopStore.GetCurrentTelopMessage()
	if telopMessage.IsSome() {
		h.wsService.PushTelop(telopMessage.Unwrap())
	}

	h.AudioService.SetIsConversion(true)
	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
}

// PostConversionCMMode godoc
// @Summary      conversion cm-mode
// @Description  endpoint for conversion to cm-mode
// @Tags         conversion
// @Accept       json
// @Produce      json
// @Param CMStateStart body requests.CMStateRequest true "post CM request"
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /conversion/cm-mode [post]
func (h *ConversionHandlers) PostConversionCMMode(c *gin.Context) {
	var conv requests.CMStateRequest
	if err := c.ShouldBindJSON(&conv); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	convEntity := conv.ToDomainCMState()

	//ここでのテロップはいらないかも
	// h.TelopStore.SetConversionTelop(convEntity)
	// telopMessage := h.TelopStore.GetCurrentTelopMessage()
	// if telopMessage.IsSome() {
	// 	h.wsService.PushTelop(telopMessage.Unwrap())
	// }

	if convEntity.IsCMMode {
		err := h.AudioService.SetCMScene()
		if err != nil {
			errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
				Kind: entities.InvalidFormat})
			c.JSON(status, errRes)
			return
		}
	} else {
		err := h.AudioService.SetMute(false)
		if err != nil {
			errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
				Kind: entities.InvalidFormat})
			c.JSON(status, errRes)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
}
