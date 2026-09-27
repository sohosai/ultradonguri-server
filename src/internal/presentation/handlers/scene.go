package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/repositories"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"
)

type SceneHandlers struct {
	SceneManager repositories.SceneManager
}

// PostBurariScene godoc
// @Summary      change to burari scene
// @Description  endpoint for change OBS scene to burari-travel
// @Tags         scene
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /burari-scene [post]
func (h *SceneHandlers) PostBurariScene(c *gin.Context) {
	results := []responses.Result{}

	err := h.SceneManager.SetBurariScene()
	if err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	results = append(results, responses.Result{
		Operation: "Burari_Scene_change",
		Success:   true,
	})

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
}

// PostCmScene godoc
// @Summary      change to cm scene
// @Description  endpoint for change OBS scene to cm
// @Tags         scene
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /cm-scene [post]
func (h *SceneHandlers) PostCmScene(c *gin.Context) {
	results := []responses.Result{}

	err := h.SceneManager.SetCMScene()
	if err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	results = append(results, responses.Result{
		Operation: "CM_Scene_change",
		Success:   true,
	})

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
}

// PostNormalScene godoc
// @Summary      change to normal scene
// @Description  endpoint for change OBS scene to normal
// @Tags         scene
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /normal-scene [post]
func (h *SceneHandlers) PostNormalScene(c *gin.Context) {
	results := []responses.Result{}

	err := h.SceneManager.SetNormalScene()
	if err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	results = append(results, responses.Result{
		Operation: "Normal_Scene_change",
		Success:   true,
	})

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
}
