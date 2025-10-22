package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andreykaipov/goobs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/sohosai/ultradonguri-server/docs"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/scene"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
	"github.com/sohosai/ultradonguri-server/internal/presentation/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title       Ultradonguri API
// @version     0.1
// @description Sohosai 2025 project telop sending API
// @BasePath    /
func main() {
	ADDR := os.Getenv("ADDRESS")
	PASS := os.Getenv("PASSWORD")
	scenes := scene.SceneNames{
		Normal: os.Getenv("NORMAL_SCENE_NAME"),
		Muted:  os.Getenv("MUTED_SCENE_NAME"),
		CM:     os.Getenv("CM_SCENE_NAME"),
	}
	CONTROLLER_ORIGINS := os.Getenv("CONTROLLER_ADDRESS")
	allowOrigins := strings.Split(CONTROLLER_ORIGINS, ",")
	SCENE_BACKUP_PATH := os.Getenv("SCENE_BACKUP_PATH")

	obsClient, err := goobs.New(ADDR, goobs.WithPassword(PASS))
	if err != nil {
		log.Fatal(err)
	}
	defer obsClient.Disconnect()

	// sceneManagerのバックアップからのrestoreを試す
	sceneManager, err := scene.RestoreSceneManager(obsClient, scenes, SCENE_BACKUP_PATH)
	if err != nil {
		log.Printf("Failed to resotre scene manager: %s", err.Error())

		// 失敗した場合は新しく作成する
		sceneManager, err = scene.NewSceneManager(obsClient, scenes, SCENE_BACKUP_PATH) // nil は本番では goobs.Client を渡す
		if err != nil {
			panic(err)
		}
	}

	wsHub := websocket.NewWebSocketHub(5)
	go wsHub.StartTelopWebsocketBroadcastWorker()

	telopManager := telop.NewTelopManager()

	h := handlers.NewHandler(sceneManager, telopManager, wsHub)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			for _, prefix := range allowOrigins {
				prefix = strings.TrimSpace(prefix)
				if strings.HasPrefix(origin, prefix) {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	h.Handle(r)

	swagger := r.Group("/swagger")
	{
		// Swagger UI を /swagger/index.html で公開
		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	fmt.Println("Application Starts!")
	r.Run("0.0.0.0:8080")
}
