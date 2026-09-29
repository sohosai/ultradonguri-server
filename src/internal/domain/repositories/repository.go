package repositories

import (
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type SceneManager interface {
	SetNormalScene() error
	SetCMScene() error
	GetCurrentScene() (string, error)
	IsCm() (bool, error)
}

type TelopManager interface {
	SetPerformanceTelop(entities.Performance)
	SetMusicTelop(entities.Music)
	SetConversionTelop(entities.ConversionPost)
	GetCurrentTelopMessage() utils.Option[entities.TelopMessage]
	IsConversion() bool
	ShouldBeMuted() bool
}
