package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"
)

type CopyRightHandler struct {
	wsService *websocket.WebSocketHub
}

// PostDisplayCopyRight godoc
// @Summary      display copyright
// @Description  endpoint to display copyright
// @Tags         display-copyright
// @Accept       json
// @Produce      json
// @Param isDisplayedCopy-right body requests.DisplayCopyrightRequest true "post display-copyright"
// @Success      200  {object}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /display-copyright [post]
func (h *CopyRightHandler) PostDisplayCopyRight(c *gin.Context) {
	var disp requests.DisplayCopyrightRequest
	results := []responses.Result{}
	if err := c.ShouldBindJSON(&disp); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	c.IndentedJSON(http.StatusOK, responses.SuccessResponse{Message: "OK", Results: results})
}
