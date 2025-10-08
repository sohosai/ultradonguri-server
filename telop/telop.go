package telop

import (
	"log/slog"

	"example.com/donguri-back/spec"
	"example.com/donguri-back/util"
)

type TelopStore struct {
	TelopType    spec.TelopType
	Performance  util.Option[spec.PerformancePost]
	Conversion   util.Option[spec.ConversionPost]
	TelopMessage util.Option[spec.TelopMessage]
}

func NewTelopStore() TelopStore {
	return TelopStore{
		TelopType:    spec.TelopTypeEmpty,
		Performance:  util.None[spec.PerformancePost](),
		Conversion:   util.None[spec.ConversionPost](),
		TelopMessage: util.None[spec.TelopMessage](),
	}
}

func (self *TelopStore) SetPerformanceTelop(telop spec.PerformancePost) {
	self.Performance = util.Some(telop)
	self.Conversion = util.None[spec.ConversionPost]()
	self.TelopType = spec.TelopTypePerformance

	slog.Info("Telop changed: ", "performance", self.Performance)
}

func (self *TelopStore) SetConversionTelop(telop spec.ConversionPost) {
	self.Performance = util.None[spec.PerformancePost]()
	self.Conversion = util.Some(telop)
	self.TelopType = spec.TelopTypeConversion

	slog.Info("Telop changed: ", "conversion", self.Conversion)
}

func (self *TelopStore) GetCurrentTelopMessage() util.Option[spec.TelopMessage] {
	switch self.TelopType {
	case spec.TelopTypePerformance:
		{
			if self.Performance.IsNone() {
				return util.None[spec.TelopMessage]()
			}

			performance := self.Performance.Unwrap()

			message := spec.TelopMessage{
				Type:            spec.TelopTypePerformance,
				PerformanceData: &performance,
				ConversionData:  nil,
			}

			return util.Some(message)

		}
	case spec.TelopTypeConversion:
		{
			if self.Conversion.IsNone() {
				return util.None[spec.TelopMessage]()
			}

			conversion := self.Conversion.Unwrap()

			message := spec.TelopMessage{
				Type:            spec.TelopTypeConversion,
				PerformanceData: nil,
				ConversionData:  &conversion,
			}

			return util.Some(message)
		}
	case spec.TelopTypeEmpty:
		{
			return util.None[spec.TelopMessage]()
		}
	}

	return util.None[spec.TelopMessage]()
}
