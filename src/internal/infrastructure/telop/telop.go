package telop

import (
	"encoding/json"
	"log"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type TelopStore struct {
	TelopType    entities.TelopType
	Performance  utils.Option[entities.PerformancePost]
	Conversion   utils.Option[entities.ConversionPost]
	TelopMessage utils.Option[entities.TelopMessage]
}

func NewTelopStore() *TelopStore {
	return &TelopStore{
		TelopType:   entities.TelopTypeEmpty,
		Performance: utils.None[entities.PerformancePost](),
		Conversion:  utils.None[entities.ConversionPost](),
	}
}

func (self *TelopStore) SetPerformanceTelop(telop entities.Performance) {
	self.Performance = utils.Some(entities.PerformancePost{
		Performance: telop,
	})
	self.Conversion = utils.None[entities.ConversionPost]()
	self.TelopType = entities.TelopTypePerformance

	telopJson, _ := json.Marshal(telop)
	// slog.Info("Telop changed: ", "performance", self.Performance)
	log.Printf("Telop changed: %s", string(telopJson))

}

func (self *TelopStore) SetMusicTelop(telop entities.Music) {
	self.Performance = utils.Some(entities.PerformancePost{
		Music: telop,
	})
	self.Conversion = utils.None[entities.ConversionPost]()
	self.TelopType = entities.TelopTypePerformance

	telopJson, _ := json.Marshal(telop)
	// slog.Info("Telop changed: ", "performance", self.Performance)
	log.Printf("Telop changed: %s", string(telopJson))

}

func (self *TelopStore) SetConversionTelop(telop entities.ConversionPost) {
	self.Performance = utils.None[entities.PerformancePost]()
	self.Conversion = utils.Some(telop)
	self.TelopType = entities.TelopTypeConversion

	telopJson, _ := json.Marshal(telop)
	// slog.Info("Telop changed: ", "conversion", self.Conversion)
	log.Printf("Telop changed: %s", string(telopJson))
}

func (self *TelopStore) GetCurrentTelopMessage() utils.Option[entities.TelopMessage] {
	switch self.TelopType {
	case entities.TelopTypePerformance:
		{
			if self.Performance.IsNone() {
				return utils.None[entities.TelopMessage]()
			}

			performance := self.Performance.Unwrap()

			message := entities.TelopMessage{
				Type:            entities.TelopTypePerformance,
				PerformanceData: &performance,
				ConversionData:  nil,
			}

			return utils.Some(message)

		}
	case entities.TelopTypeConversion:
		{
			if self.Conversion.IsNone() {
				return utils.None[entities.TelopMessage]()
			}

			conversion := self.Conversion.Unwrap()

			message := entities.TelopMessage{
				Type:            entities.TelopTypeConversion,
				PerformanceData: nil,
				ConversionData:  &conversion,
			}

			return utils.Some(message)
		}
	case entities.TelopTypeEmpty:
		{
			return utils.None[entities.TelopMessage]()
		}
	}

	return utils.None[entities.TelopMessage]()
}
