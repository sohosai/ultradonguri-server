package repositories

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type AudioService interface {
	Mute() error
	UnMute() error
	SetMute(bool) error
	GetMute() (bool, error)
}

type TelopService interface {
	SetPerformanceTelop(entities.PerformancePost)
	SetConversionTelop(entities.ConversionPost)
}
