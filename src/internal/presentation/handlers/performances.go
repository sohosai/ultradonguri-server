package handlers

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/presentation/model/responses"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/file"
)

type PerformancesHandler struct{}

func (h *PerformancesHandler) GetPerformances(c *gin.Context) {
	perfs, err := file.GetPerformances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newPerfState := responses.NewPerformancesResponse(perfs)
	c.JSON(http.StatusOK, newPerfState)
}
