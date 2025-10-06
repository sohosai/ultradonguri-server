package telop

import (
	"log/slog"

	"example.com/donguri-back/spec"
	"example.com/donguri-back/util"
)

type TelopType string

const (
	Performance TelopType = "performance"
	Conversion  TelopType = "conversion"
	Empty       TelopType = "empty"
)

type TelopClient struct {
	TelopType   TelopType
	Performance util.Option[spec.PerformancePost]
	Conversion  util.Option[spec.ConversionPost]
}

func NewTelopClient() TelopClient {
	return TelopClient{
		TelopType:   Empty,
		Performance: util.None[spec.PerformancePost](),
		Conversion:  util.None[spec.ConversionPost](),
	}
}

func (self *TelopClient) SetPerformanceTelop(telop spec.PerformancePost) {
	self.Performance = util.Some(telop)
	self.Conversion = util.None[spec.ConversionPost]()
	self.TelopType = Performance

	slog.Info("Telop changed: ", "performance", self.Performance)
}

func (self *TelopClient) SetConversionTelop(telop spec.ConversionPost) {
	self.Performance = util.None[spec.PerformancePost]()
	self.Conversion = util.Some(telop)
	self.TelopType = Conversion

	slog.Info("Telop changed: ", "conversion", self.Conversion)
}
