package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type ConversionRequest struct {
	NextPerformances []NextPerformanceRequest `json:"next_performances"`
}

type NextPerformanceRequest struct {
	Title       string      `json:"title"`
	Performer   string      `json:"performer"`
	Description string      `json:"description"`
	StartsAt    entities.HM `json:"starts_at"`
}

type CMStateRequest struct {
	IsCMMode bool `json:"is_cm_mode"`
}

func (conv ConversionRequest) ToDomainConversion() entities.ConversionPost {
	var nextPerformances []entities.NextPerformance
	for _, e := range conv.NextPerformances {
		nextPerformances = append(nextPerformances, entities.NextPerformance{
			Title:       e.Title,
			Performer:   e.Performer,
			Description: e.Description,
			StartsAt:    e.StartsAt,
		})
	}
	return entities.ConversionPost{
		NextPerformances: nextPerformances,
	}
}

func (cm CMStateRequest) ToDomainCMState() entities.CMState {
	return entities.CMState{
		IsCMMode: cm.IsCMMode,
	}
}
