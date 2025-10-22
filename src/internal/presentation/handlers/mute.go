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
	results := []responses.Result{}
	if err := c.ShouldBindJSON(&muteReq); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	newMuteState := muteReq.ToDomainMute() //domainの型に変換

	if newMuteState.IsMuted { // 強制ミュートをする場合

		// CM中の場合はMuteに切り替えない
		isCm, err := h.SceneManager.IsCm()
		if err != nil {
			errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
				Kind: entities.InvalidFormat})
			c.JSON(status, errRes)
		}
		if isCm {
			errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: "cannot force_mute when CM scene",
				Kind: entities.CannotChangeState})
			c.JSON(status, errRes)
			return
		}

		// SceneをMuteに切り替える
		h.SceneManager.SetForceMuteFlag(true)
		results = append(results, responses.Result{
			Operation: "force_mute_on",
			Success:   true,
		})
		err = h.SceneManager.SetMute(true)
		results = append(results, responses.Result{
			Operation: "mute_change",
			Success:   err == nil,
		})

		c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
		return
	}

	// 強制ミュートを解除する場合

	isCm, err := h.SceneManager.IsCm()
	if err != nil {
		// エラーは仮
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.CannotForceMute})
		c.JSON(status, errRes)
		return
	}

	// forceMuteFlagを無効化する
	h.SceneManager.SetForceMuteFlag(false)
	results = append(results, responses.Result{
		Operation: "force_mute_off",
		Success:   true,
	})

	if (h.TelopManager.IsConversion() || !h.TelopManager.ShouldBeMuted()) && !isCm {
		// 現在のTelopがConversion
		// または
		// 現在のTelopがPerformanceでMusicがshould_be_muted=falseの場合
		// かつCMモードじゃない場合

		// SceneをNormalへ移行する
		err := h.SceneManager.SetNormalScene()
		results = append(results, responses.Result{
			Operation: "mute_change",
			Success:   err == nil,
		})

		c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
		return
	}

	// ミュート状態自体は継続する場合

	// 結果をエラーにするのかどうかは後で決める
	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
}
