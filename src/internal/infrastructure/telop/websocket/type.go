package websocket

import "reflect"

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

var typeRegistry = map[reflect.Type]string{
	reflect.TypeOf(PerformanceStartData{}): TypePerformanceStart,
	reflect.TypeOf(PerformanceMusicData{}): TypePerformanceMusic,
	reflect.TypeOf(PerformanceMusicData{}): TypeConversionStart,
	reflect.TypeOf(ConversionCmMode{}):     TypeConversionCmMode,
}
