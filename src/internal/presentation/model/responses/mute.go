package responses

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

//requestのほうと同じなのでまとめることもできる。
type MuteStateResponse struct {
	IsMuted bool `json:"is_muted"`
}

func NewMuteStateResponse(muteState entities.MuteState) MuteStateResponse {
	return MuteStateResponse{
		IsMuted: muteState.IsMuted,
	}
}
