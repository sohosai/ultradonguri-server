package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

//responsesのほうと同じなのでまとめることもできる。
type MuteStateRequest struct {
	IsMuted bool `json:"is_muted"`
}

func NewMuteStateRequest(muteState entities.MuteState) MuteStateRequest {
	return MuteStateRequest{
		IsMuted: muteState.IsMuted,
	}
}
