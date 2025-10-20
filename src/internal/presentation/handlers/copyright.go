package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/requests"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"
)

type CopyRightHandler struct {
	wsService *websocket.WebSocketHub
}

// HealthCheck godoc
// @Summary      desplay copyright
// @Description  endpoint to desplay copyright
// @Tags         desplay-copyright
// @Produce      json
// @Success      200  {oboject}  responses.SuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /desplay-copyright [post]
func (h *CopyRightHandler) PostDesplayCopyRight(c *gin.Context) {
	var desp requests.DesplayCopyrightRequest
	if err := c.ShouldBindJSON(&desp); err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidFormat})
		c.JSON(status, errRes)
		return
	}

	despEntity := desp.ToDomainCopyright()

	resp, err := websocket.TypedWebSocketResponse[websocket.DesplayCopyrightData]{
		Type: websocket.TypeConversionStart,
		Data: websocket.ToDataDesplayCopyright(despEntity),
	}.Encode()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.wsService.PushTelop(resp)

	c.IndentedJSON(http.StatusOK, responses.SuccessResponse{Message: "OK"})
}
