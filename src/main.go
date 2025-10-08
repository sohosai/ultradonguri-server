package main

import (
	"fmt"
	"log"

	"github.com/andreykaipov/goobs"
	"github.com/gin-gonic/gin"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/audio"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/file"
	"github.com/sohosai/ultradonguri-server/internal/infrastructure/telop"
	"github.com/sohosai/ultradonguri-server/internal/presentation/handlers"
)

func main() {
	ADDR := "192.168.0.38:4455"
	PASS := "FiXyBLyCug7xt2Zw"
	SCENE := "Dev"
	AUDIO_INPUT := "Audio Input"

	// OBS との接続を確立
	obsClient, err := goobs.New(ADDR, goobs.WithPassword(PASS))
	if err != nil {
		log.Fatal(err)
	}
	defer obsClient.Disconnect()

	// OBS AudioClient 初期化
	audioClient, err := audio.NewAudioClient(obsClient, SCENE, AUDIO_INPUT) // nil は本番では goobs.Client を渡す
	if err != nil {
		panic(err)
	}

	// audioClient := obsClient

	telopClient := telop.NewTelopClient()

	perfRepo := &file.PerformanceRepository{Path: "events.json"}

	h := handlers.NewHandler(audioClient, telopClient, perfRepo)

	r := gin.Default()
	h.Handle(r)

	fmt.Println("Application Starts!")
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
