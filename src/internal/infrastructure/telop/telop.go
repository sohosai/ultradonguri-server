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

func NewTelopClient() *TelopStore {
	return &TelopStore{
		TelopType:   entities.TelopTypeEmpty,
		Performance: utils.None[entities.PerformancePost](),
		Conversion:  utils.None[entities.ConversionPost](),
	}
}

func (self *TelopStore) SetPerformanceTelop(telop entities.PerformancePost) {
	self.Performance = utils.Some(telop)
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
