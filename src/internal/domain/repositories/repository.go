package repositories

type SceneManager interface {
	SetNormalScene() error
	SetCMScene() error
	SetBurariScene() error
	GetCurrentScene() (string, error)
	IsCm() (bool, error)
}
