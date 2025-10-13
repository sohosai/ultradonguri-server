package repositories

import (
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type AudioService interface {
	Mute() error
	UnMute() error
	SetMute(bool) error
	GetMute() (entities.MuteState, error)
}

type TelopStore interface {
	SetPerformanceTelop(entities.PerformancePost)
	SetConversionTelop(entities.ConversionPost)
	GetCurrentTelopMessage() utils.Option[entities.TelopMessage]
}
