package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"

	"github.com/gin-gonic/gin"
)

type MuteHandler struct {
	SceneManager repositories.SceneManager
	TelopManager repositories.TelopManager
}

// PostForceMute godoc
// @Summary      force mute
// @Description  endpoint for force mute and prevent changes
// @Tags         force_mute
// @Accept       json
// @Produce      json
// @Param isMuted body requests.MuteStateRequest true "select mute state"
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /force_mute [post]
func (h *MuteHandler) PostForceMuted(c *gin.Context) {
	var muteReq requests.MuteStateRequest //jsonを受け取るため
	if err := c.ShouldBindJSON(&muteReq); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	newMuteState := muteReq.ToDomainMute() //domainの型に変換

	if newMuteState.IsMuted { // 強制ミュートをする場合
		// forceMuteFlagを有効化する
		// h.SceneManager.SetForceMuteFlag(true)

		// CM中の場合はMuteに切り替えない
		isCm, err := h.SceneManager.IsCm()
		if err != nil {
			errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
				Kind: entities.InvalidFormat})
			c.JSON(status, errRes)
		}
		if isCm {
			// c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
			errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: "cannot force_mute when CM scene",
				Kind: entities.CannotChangeState})
			c.JSON(status, errRes)
			return
		}

		// SceneをMuteに切り替える
		h.SceneManager.SetForceMuteFlag(true)
		if err := h.SceneManager.SetMute(true); err != nil {
			// エラーは仮
			errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
				Kind: entities.CannotForceMute})
			c.JSON(status, errRes)
			return
		}

		c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
		return
	}

	// 強制ミュートを解除する場合

	// forceMuteFlagを無効化する
	h.SceneManager.SetForceMuteFlag(false)

	isCm, err := h.SceneManager.IsCm()
	if err != nil {
		// エラーは仮
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.CannotForceMute})
		c.JSON(status, errRes)
		return
	}

	if (h.TelopManager.IsConversion() || !h.TelopManager.ShouldBeMuted()) && !isCm {
		// 現在のTelopがConversion
		// または
		// 現在のTelopがPerformanceでMusicがshould_be_muted=falseの場合
		// かつCMモードじゃない場合

		// SceneをNormalへ移行する
		h.SceneManager.SetNormalScene()

		c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
		return
	}

	// ミュート状態自体は継続する場合

	// 結果をエラーにするのかどうかは後で決める
	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
}
