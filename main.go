package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"example.com/donguri-back/audio.go"
	"example.com/donguri-back/client"
	"example.com/donguri-back/event"
	"github.com/andreykaipov/goobs"
	"github.com/gin-gonic/gin"
)

type IdReq struct {
	Id string `json:"id"`
}

type MuteReq struct {
	Mute bool `json:"mute"`
}

func main() {
	ADDR := "127.0.0.1:4455"
	PASS := "obs-websocket-password"
	SCENE := "Dev"
	AUDIO_INPUT := "Audio Input"

	// OBS との接続を確立
	obsClient, err := goobs.New(ADDR, goobs.WithPassword(PASS))
	if err != nil {
		log.Fatal(err)
	}
	defer obsClient.Disconnect()

	sharedClient := client.NewSharedClient(obsClient)

	// 曲名テロップのTelopClientを作成
	musicTelop, err := event.NewMusicTelopClient(sharedClient, SCENE)
	if err != nil {
		fmt.Printf("Failed to create music telop client: %s\n", err)
	}

	// パフォーマンス名テロップのTelopClientを作成
	eventTelop, err := event.NewEventTelopClient(sharedClient, SCENE)
	if err != nil {
		fmt.Printf("Failed to create event telop client: %s\n", err)
	}

	// AudioClientの作成
	audioClient, err := audio.NewAudioClient(sharedClient, SCENE, AUDIO_INPUT)
	if err != nil {
		fmt.Printf("Failed to create audio client: %s\n", err)
	}

	// webサーバー
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// music api
	// 現在の楽曲を取得する
	r.GET("/v1/music", func(ctx *gin.Context) {
		music, err := musicTelop.GetMusic()

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get current music: " + err.Error()})
			return
		}

		ctx.JSON(200, music)
	})

	// 楽曲を設定する
	r.POST("/v1/music", func(ctx *gin.Context) {
		var body IdReq
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		music, err := event.FindMusicById(body.Id)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find the music: " + err.Error()})
			return
		}

		musicTelop.SetMusic(music)
		musicTelop.ApplyMusicChange()
		if music.ShouldBeMuted {
			audioClient.Mute()
		} else {
			audioClient.UnMute()
		}

		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// event api
	// イベント一覧を取得する
	r.GET("/v1/events", func(ctx *gin.Context) {
		events, err := event.GetEvents()

		if err != nil {
			slog.Error("Failed to get events", "err", err)
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, events)
	})

	// 現在のイベントを取得する
	r.GET("/v1/event", func(ctx *gin.Context) {
		event, err := eventTelop.GetEvent()

		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}

		ctx.JSON(200, event)
	})

	// イベントを設定する
	r.POST("/v1/event", func(ctx *gin.Context) {
		var body IdReq
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		event, err := event.FindEventById(body.Id)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to find the event: " + err.Error()})
			return
		}

		eventTelop.SetEvent(event)
		eventTelop.ApplyEventChange()

		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// AudioのAPI
	r.GET("/v1/mute", func(ctx *gin.Context) {
		isMuted, err := audioClient.GetMute()

		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to get state of mute: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"muted": isMuted})
	})

	r.POST("/v1/mute", func(ctx *gin.Context) {
		var muteReq MuteReq
		if err := ctx.ShouldBindJSON(&muteReq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		err := audioClient.SetMute(muteReq.Mute)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set mute state: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// デフォルト :8080 で起動
	_ = r.Run(":8080")
}
