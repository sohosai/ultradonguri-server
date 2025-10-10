package telop

import (
	"encoding/json"
	"log"
	"log/slog"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/domain/entities/telopentity"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type TelopType string

type TelopClient struct {
	TelopType   telopentity.TelopType
	Performance utils.Option[entities.PerformancePost]
	Conversion  utils.Option[entities.ConversionPost]
}

func NewTelopClient() *TelopClient {
	return &TelopClient{
		TelopType:   telopentity.Empty,
		Performance: utils.None[entities.PerformancePost](),
		Conversion:  utils.None[entities.ConversionPost](),
	}
}

func (self *TelopClient) SetPerformanceTelop(telop entities.PerformancePost) {
	self.Performance = utils.Some(telop)
	self.Conversion = utils.None[entities.ConversionPost]()
	self.TelopType = telopentity.Performance

	telopJson, _ := json.Marshal(telop)
	// slog.Info("Telop changed: ", "performance", self.Performance)
	log.Printf("Telop changed: %s", string(telopJson))

}

func (self *TelopClient) SetConversionTelop(telop entities.ConversionPost) {
	self.Performance = utils.None[entities.PerformancePost]()
	self.Conversion = utils.Some(telop)
	self.TelopType = telopentity.Conversion

	slog.Info("Telop changed: ", "conversion", self.Conversion)
}
