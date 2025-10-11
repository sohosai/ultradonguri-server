package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type ConversionRequest struct {
	NextPerformances []NextPerformanceRequest `json:"next_performances"`
}

type NextPerformanceRequest struct {
	Title       string `json:"title"`
	Performer   string `json:"performer"`
	Description string `json:"description"`
}

func NewConversionRequest(entities entities.ConversionPost) ConversionRequest {
	var nextPerformances []NextPerformanceRequest
	for _, e := range entities.NextPerformances {
		nextPerformances = append(nextPerformances, NextPerformanceRequest{
			Title:       e.Title,
			Performer:   e.Performer,
			Description: e.Description,
		})
	}
	return ConversionRequest{
		NextPerformances: nextPerformances,
	}
}
