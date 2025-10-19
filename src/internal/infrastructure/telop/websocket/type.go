package websocket

import (
	"reflect"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
)

type PerformanceStartData struct {
	Title     string `json:"title"`
	Performer string `json:"performer"`
}

type PerformanceMusicData struct {
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	ShouldBeMuted bool   `json:"should_be_muted"`
}

type ConversionStart struct {
	NextPerformances []NextPerformanceWs `json:"next_performances"`
}

type ConversionCmMode struct {
	IsCMMode bool `json:"is_cm_mode"`
}

type NextPerformanceWs struct {
	Title       string `json:"title"`
	Performer   string `json:"performer"`
	Description string `json:"description"`
}

const (
	TypePerformanceStart = "/performance/start"
	TypePerformanceMusic = "/performance/music"
	TypeConversionStart  = "/conversion/start"
	TypeConversionCmMode = "/conversion/cm-mode"
)

// dataに乗せる型とエンドポイントが正しいかの判断用
var typeRegistry = map[reflect.Type]string{
	reflect.TypeOf(PerformanceStartData{}): TypePerformanceStart,
	// reflect.TypeOf(entities.Performance{}): TypePerformanceStart,
	reflect.TypeOf(PerformanceMusicData{}): TypePerformanceMusic,
	reflect.TypeOf(ConversionStart{}):      TypeConversionStart,
	reflect.TypeOf(ConversionCmMode{}):     TypeConversionCmMode,
}

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

func ToDataConvStart(c entities.ConversionPost) ConversionStart {
	nextPerformances := make([]NextPerformanceWs, len(c.NextPerformances))
	for i, np := range c.NextPerformances {
		nextPerformances[i] = NextPerformanceWs{
			Title:       np.Title,
			Performer:   np.Performer,
			Description: np.Description,
		}
	}

	return ConversionStart{
		NextPerformances: nextPerformances,
	}
}

func ToDataConvCmMode(c entities.CMState) ConversionCmMode {
	return ConversionCmMode{
		IsCMMode: c.IsCMMode,
	}
}
