package repositories

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
