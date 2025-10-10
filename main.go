package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"example.com/donguri-back/audio"
	"example.com/donguri-back/spec"
	"example.com/donguri-back/telop"
	"github.com/andreykaipov/goobs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	telopStore := telop.NewTelopStore()

	// AudioClientの作成
	audioClient, err := audio.NewAudioClient(obsClient, SCENE, AUDIO_INPUT)
	if err != nil {
		fmt.Printf("Failed to create audio client: %s\n", err)
	}

	// Telopをwebsocketにブロードキャストするためのworkerを起動
	go StartTelopWebsocketBroadcastWorker()

	// websocket upgrader
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return false
			}
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			switch u.Hostname() {
			case "localhost", "127.0.0.1", "::1":
				return true
			default:
				return false
			}
		},
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

		telopStore.SetPerformanceTelop(performancePost)
		telopMessage := telopStore.GetCurrentTelopMessage()
		if telopMessage.IsSome() {
			PushTelop(telopMessage.Unwrap())
		}

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

		telopStore.SetConversionTelop(conversionPost)
		telopMessage := telopStore.GetCurrentTelopMessage()
		if telopMessage.IsSome() {
			PushTelop(telopMessage.Unwrap())
		}

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

	r.GET("/ws", func(c *gin.Context) {
		wsConnection, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		slog.Info("New websocket connection is established")

		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		WsConnections.Lock()
		WsConnections.conns[wsConnection] = true
		for conn := range WsConnections.conns {
			slog.Info("", "Addr: ", conn.LocalAddr().String())
		}
		WsConnections.Unlock()

		defer func() {
			WsConnections.Lock()
			delete(WsConnections.conns, wsConnection)
			WsConnections.Unlock()
			wsConnection.Close()
		}()

		for {
			if _, _, err := wsConnection.ReadMessage(); err != nil {
				slog.Error("read: ", err)
				break
			}
		}

	})

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

	if dec.More() {
		return nil, fmt.Errorf("trailing data after JSON array")
	}
	return events, nil
}
