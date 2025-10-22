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
	TelopManager repositories.TelopManager
	SceneManager repositories.SceneManager
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
	results := []responses.Result{}
	if err := c.ShouldBindJSON(&conv); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	convEntity := conv.ToDomainConversion()

	// TelopをConversionへ切り替え
	h.TelopManager.SetConversionTelop(convEntity)

	// viewerへの通知
	resp, err := websocket.TypedWebSocketResponse[websocket.ConversionStartData]{
		Type: websocket.TypeConversionStart,
		Data: websocket.ToDataConvStart(convEntity),
	}.Encode()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.wsService.PushTelop(resp)
	results = append(results, responses.Result{
		Operation: "telop_change",
		Success:   true,
	})

	// Normalシーンへ切り替え
	err = h.SceneManager.SetNormalScene()
	results = append(results, responses.Result{
		Operation: "CM_Scene_change",
		Success:   err == nil,
	})

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
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
	results := []responses.Result{}
	if err := c.ShouldBindJSON(&conv); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	convEntity := conv.ToDomainCMState()

	if h.TelopManager.IsConversion() {
		// 転換パートでのみViewerへの通知とシーンの切り替えを行う

		// シーンの切り替え
		if convEntity.IsCMMode { // CMシーンへの切り替えを指定された場合
			// シーンをCMに切り替える
			err := h.SceneManager.SetCMScene()
			if err != nil {
				// エラーは仮
				errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
					Kind: entities.InvalidFormat})
				c.JSON(status, errRes)
				return
			} else {
				results = append(results, responses.Result{
					Operation: "CM_Scene_change",
					Success:   true,
				})
			}
		} else { // CMシーンからNormalへ戻る場合
			// CMシーンに切り替わるのはConversion中だけで、
			// 切り替えの際にTelopの情報は消されずに維持されるのでシーンだけNormalに戻せば良い
			if err := h.SceneManager.SetNormalScene(); err != nil {
				errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
					Kind: entities.InvalidFormat})
				c.JSON(status, errRes)
				return
			} else {
				results = append(results, responses.Result{
					Operation: "mute_change",
					Success:   true,
				})
			}
		}

		// Viewerへの通知
		resp, err := websocket.TypedWebSocketResponse[websocket.ConversionCmModeData]{
			Type: websocket.TypeConversionCmMode,
			Data: websocket.ToDataConvCmMode(convEntity),
		}.Encode()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.wsService.PushTelop(resp)
		results = append(results, responses.Result{
			Operation: "telop_change",
			Success:   true,
		})

		c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
		return
	}

	// 転換パートでない場合はエラー
	// エラー処理は仮
	errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: "転換パートじゃないですよ的なエラーを発生させたい",
		Kind: entities.InvalidFormat})
	c.JSON(status, errRes)
}
