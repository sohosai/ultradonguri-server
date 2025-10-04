package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/presentation/handlers"
)

func main() {
	r := gin.Default()
	handlers.Handle(r)

	fmt.Println("Application Starts!")
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
