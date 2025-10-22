package websocket

import (
	"reflect"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
)

// ws送信用の型
type PerformanceStartData struct {
	Title     string `json:"title"`
	Performer string `json:"performer"`
}

type PerformanceMusicData struct {
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	ShouldBeMuted bool   `json:"should_be_muted"`
}

type ConversionStartData struct {
	NextPerformances []NextPerformanceData `json:"next_performances"`
}

type ConversionCmModeData struct {
	IsCMMode bool `json:"is_cm_mode"`
}

type NextPerformanceData struct {
	Title       string      `json:"title"`
	Performer   string      `json:"performer"`
	Description string      `json:"description"`
	StartsAt    entities.HM `json:"starts_at"`
}

type DisplayCopyrightData struct {
	IsDisplay bool `json:"is_displayed_copyright"`
}

type WsMessageType string

const (
	TypePerformanceStart WsMessageType = "/performance/start"
	TypePerformanceMusic WsMessageType = "/performance/music"
	TypeConversionStart  WsMessageType = "/conversion/start"
	TypeConversionCmMode WsMessageType = "/conversion/cm-mode"
	TypeDisplayCopyright WsMessageType = "/display-copyright"
)

// wsで送信するdataに乗せる型とエンドポイントが正しいかの判断用のマップ
var typeRegistry = map[reflect.Type]WsMessageType{
	reflect.TypeOf(PerformanceStartData{}): TypePerformanceStart,
	reflect.TypeOf(PerformanceMusicData{}): TypePerformanceMusic,
	reflect.TypeOf(ConversionStartData{}):  TypeConversionStart,
	reflect.TypeOf(ConversionCmModeData{}): TypeConversionCmMode,
	reflect.TypeOf(DisplayCopyrightData{}): TypeDisplayCopyright,
}

// ドメイン層の型にjsonのフィールドをつけたws送信用の型に変換する関数
func ToDataPerfStart(p entities.Performance) PerformanceStartData {
	return PerformanceStartData{
		Title:     p.Title,
		Performer: p.Performer,
	}
}

func ToDataPerfMusic(p entities.Music) PerformanceMusicData {
	return PerformanceMusicData{
		Title:         p.Title,
		Artist:        p.Artist,
		ShouldBeMuted: p.ShouldBeMuted,
	}
}

func ToDataConvStart(c entities.ConversionPost) ConversionStartData {
	nextPerformances := make([]NextPerformanceData, len(c.NextPerformances))
	for i, np := range c.NextPerformances {
		nextPerformances[i] = NextPerformanceData{
			Title:       np.Title,
			Performer:   np.Performer,
			Description: np.Description,
			StartsAt:    np.StartsAt,
		}
	}

	return ConversionStartData{
		NextPerformances: nextPerformances,
	}
}

func ToDataConvCmMode(c entities.CMState) ConversionCmModeData {
	return ConversionCmModeData{
		IsCMMode: c.IsCMMode,
	}
}

func ToDataDisplayCopyright(d entities.DisplayCopyright) DisplayCopyrightData {
	return DisplayCopyrightData{
		IsDisplay: d.IsDisplay,
	}
}
