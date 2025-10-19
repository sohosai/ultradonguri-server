package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	message := "Hello World"
	c.IndentedJSON(http.StatusOK, message)
}
