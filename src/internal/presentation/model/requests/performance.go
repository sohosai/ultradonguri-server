package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type PerformanceRequest struct {
	Title     string  `json:"title"`
	Performer *string `json:"performer"`
}

func NewPerformanceRequestToEntity(p entities.Performance) PerformanceRequest {
	return PerformanceRequest{
		Title:     p.Title,
		Performer: p.Performer,
	}
}
