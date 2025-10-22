package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

//responsesのほうと同じなのでまとめることもできる。
type MuteStateRequest struct {
	// binding:"required"の影響でfalseを渡したときにエラーが出るので、bool型ではなくポインタを渡す
	IsMuted *bool `json:"is_muted" binding:"required"`
}

func (m MuteStateRequest) ToDomainMute() entities.MuteState {
	return entities.MuteState{
		IsMuted: *m.IsMuted,
	}
}
