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

			// force_mute中はmutedに移行する
			if h.SceneManager.IsForceMutedFlag() {
				if err := h.SceneManager.SetMutedScene(); err != nil {
					errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
						Kind: entities.InvalidFormat})
					c.JSON(status, errRes)
					return
				}

				results = append(results, responses.Result{
					Operation: "mute_change",
					Success:   true,
				})
			} else {
				if err := h.SceneManager.SetNormalScene(); err != nil {
					errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
						Kind: entities.InvalidFormat})
					c.JSON(status, errRes)
					return
				} else {
					results = append(results, responses.Result{
						Operation: "Normal_Scene_change",
						Success:   true,
					})
				}
			}
		}

		c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
		return
	}

	// 転換パートでない場合はエラー
	// エラー処理は仮
	errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: "転換パートじゃないですよ的なエラーを発生させたい",
		Kind: entities.InvalidFormat})
	c.JSON(status, errRes)
}
