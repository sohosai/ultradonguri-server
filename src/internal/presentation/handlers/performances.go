package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/file"
)

type PerformancesHandler struct{}

// GetPerformances godoc
// @Summary      get performances data
// @Description  endpoint for get performances data
// @Tags         performances
// @Produce      json
// @Success      200  {object}  []responses.PerformancesResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /performances [get]
func (h *PerformancesHandler) GetPerformances(c *gin.Context) {
	perfs, err := file.GetPerformances()
	if err != nil {
		errRes, status := responses.NewErrorResponseAndHTTPStatus(entities.AppError{Message: err.Error(),
			Kind: entities.InvalidPerformancesJson})
		c.JSON(status, errRes)
		return
	}

	newPerfState := responses.NewPerformancesResponse(perfs)
	c.JSON(http.StatusOK, newPerfState)
}
