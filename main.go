package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/donguri-back/audio"
	"example.com/donguri-back/spec"
	"example.com/donguri-back/telop"
	"example.com/donguri-back/util"
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

	sharedClient := util.NewMutex(obsClient)

	telopClient := telop.NewTelopClient()

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

	r.GET("/performances", func(c *gin.Context) {
		performances, err := GetPerformances()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get performances: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, performances)
	})

	r.POST("/performance", func(c *gin.Context) {
		var performancePost spec.PerformancePost
		if err := c.ShouldBindJSON(&performancePost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		telopClient.SetPerformanceTelop(performancePost)

		if performancePost.Music.ShouldBeMuted {
			audioClient.Mute()
		} else {
			audioClient.UnMute()
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.POST("/conversion", func(c *gin.Context) {
		var conversionPost spec.ConversionPost
		if err := c.ShouldBindJSON(&conversionPost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		telopClient.SetConversionTelop(conversionPost)

		audioClient.UnMute()

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.POST("/force_mute", func(c *gin.Context) {
		var muteState spec.MuteState
		if err := c.ShouldBindJSON(&muteState); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		if err = audioClient.SetMute(muteState.IsMuted); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.GET("/force_mute", func(c *gin.Context) {
		isMuted, err := audioClient.GetMute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get mute state: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, spec.MuteState{IsMuted: isMuted})
	})

	// デフォルト :8080 で起動
	_ = r.Run(":8080")
}

func GetPerformances() ([]spec.PerformanceForPerformances, error) {
	file, err := os.Open("events.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []spec.PerformanceForPerformances
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&events); err != nil {
		return nil, err
	}

	// 余剰トークン検出
	if dec.More() {
		return nil, fmt.Errorf("trailing data after JSON array")
	}
	return events, nil
}
