package repositories

type SceneManager interface {
	SetNormalScene() error
	SetCMScene() error
	GetCurrentScene() (string, error)
	IsCm() (bool, error)
}
