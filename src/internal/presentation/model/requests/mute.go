package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

//responsesのほうと同じなのでまとめることもできる。
type MuteStateRequest struct {
	IsMuted bool `json:"is_muted"`
}

func (m MuteStateRequest) ToDomainMute() entities.MuteState {
	return entities.MuteState{
		IsMuted: m.IsMuted,
	}
}
