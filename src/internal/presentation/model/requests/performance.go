package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type PerformancePostRequest struct {
	Music       MusicRequest       `json:"music"`
	Performance PerformanceRequest `json:"performance"`
}

type MusicRequest struct {
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	ShouldBeMuted bool   `json:"should_be_muted"`
}

type PerformanceRequest struct {
	Title     string  `json:"title"`
	Performer *string `json:"performer"` // null を許すためポインタ
}

//現在(10/13)の仕様では必要になりそうなので残しておく
// type PerformanceRequest struct {
// 	Title     string  `json:"title"`
// 	Performer *string `json:"performer"`
// }

// func (p PerformanceRequest) ToDomainPerformance() entities.Performance {
// 	return entities.Performance{
// 		Title:     p.Title,
// 		Performer: p.Performer,
// 	}
// }

func (pp PerformancePostRequest) ToDomainPerformancePost() entities.PerformancePost {
	return entities.PerformancePost{
		Music: entities.Music{
			Title:         pp.Music.Title,
			Artist:        pp.Music.Artist,
			ShouldBeMuted: pp.Music.ShouldBeMuted,
		},
		Performance: entities.Performance{
			Title:     pp.Performance.Title,
			Performer: pp.Performance.Performer,
		},
	}
}
