package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type DesplayCopyrightRequest struct {
	IsDesplay bool `json:"is_desplayed_copyright"`
}

func (d DesplayCopyrightRequest) ToDomainCopyright() entities.DesplayCopyright {
	return entities.DesplayCopyright{
		IsDesplay: d.IsDesplay,
	}
}
