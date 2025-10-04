package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handle(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		message := "Hallo World"
		c.IndentedJSON(http.StatusOK, message)
	})

}
