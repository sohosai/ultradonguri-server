package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type ConversionRequest struct {
	NextPerformances []NextPerformanceRequest `json:"next_performances" binding:"required"`
}

type NextPerformanceRequest struct {
	Title       string      `json:"title" binding:"required"`
	Performer   string      `json:"performer" binding:"required"`
	Description string      `json:"description" binding:"required"`
	StartsAt    entities.HM `json:"starts_at" binding:"required"`
}

type CMStateRequest struct {
	// binding:"required"の影響でfalseを渡したときにエラーが出るので、bool型ではなくポインタを渡す
	IsCMMode *bool `json:"is_cm_mode" binding:"required"`
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
		IsCMMode: *cm.IsCMMode,
	}
}
