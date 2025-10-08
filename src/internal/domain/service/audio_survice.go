package service

// AudioService はドメインサービスとして抽象化
type AudioService interface {
	Mute() error
	UnMute() error
	SetMute(bool) error
	GetMute() (bool, error)
}
