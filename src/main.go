package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/gin-gonic/gin"

	"github.com/sohosai/ultradonguri-server/internal/infrastructure/audio"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop/websocket"
	"github.com/sohosai/ultradonguri-server/internal/presentation/handlers"
)

func main() {
	ADDR := os.Getenv("ADDRESS")
	PASS := os.Getenv("PASSWORD")
	scenes := audio.Scenes{
		Normal: os.Getenv("NORMAL_SCENE_NAME"),
		Muted:  os.Getenv("MUTED_SCENE_NAME"),
		CM:     os.Getenv("CM_SCENE_NAME"),
	}

	// fmt.Printf("ADDR: %s\n", ADDR)

	obsClient, err := goobs.New(ADDR, goobs.WithPassword(PASS))
	if err != nil {
		log.Fatal(err)
	}
	defer obsClient.Disconnect()

	audioClient, err := audio.NewAudioClient(obsClient, scenes) // nil は本番では goobs.Client を渡す
	if err != nil {
		panic(err)
	}

	wsHub := websocket.NewWebSocketHub(5)
	go wsHub.StartTelopWebsocketBroadcastWorker()

	telopStore := telop.NewTelopStore()

	h := handlers.NewHandler(audioClient, telopStore, wsHub)

	r := gin.Default()
	h.Handle(r)

	fmt.Println("Application Starts!")
	r.Run("0.0.0.0:8080")
}
