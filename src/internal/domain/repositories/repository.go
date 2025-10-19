package repositories

import (
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type AudioService interface {
	SetMute(bool) error
	// GetMute() (entities.MuteState, error)
	SetForceMute(bool) error
	SetShouldBeMuted(bool) error
	SetIsConversion(bool) error
	SetNormalScene() error
	SetMutedScene() error
	SetCMScene() error
	GetCurrentScene() (string, error)
}

type TelopStore interface {
	SetPerformanceTelop(entities.PerformancePost)
	SetConversionTelop(entities.ConversionPost)
	GetCurrentTelopMessage() utils.Option[entities.TelopMessage]
}
