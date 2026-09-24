package repositories

import (
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type SceneManager interface {
	SetMute(bool) error
	SetNormalScene() error
	SetMutedScene() error
	SetCMScene() error
	GetCurrentScene() (string, error)
	SetForceMuteFlag(bool)
	IsCm() (bool, error)
	IsForceMutedFlag() bool
}

