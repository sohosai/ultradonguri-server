package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"
)

type HealthHandler struct{}

// HealthCheck godoc
// @Summary      health check
// @Description  endpoint for health check
// @Tags         health
// @Produce      json
// @Success      200  {object}  responses.SuccessResponse
// @Router       /health [get]
func (h *HealthHandler) GetHealth(c *gin.Context) {
	message := "Hello World"
	c.IndentedJSON(http.StatusOK, responses.SuccessResponse{Message: message})
}
