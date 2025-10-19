package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type PerformancePostRequest struct {
	Music       MusicRequest       `json:"music"`
	Performance PerformanceRequest `json:"performance"`
}

type MusicPostRequest struct {
	Music MusicRequest `json:"music"`
}

type OnlyPerformancePostRequest struct {
	Performance PerformanceRequest `json:"performance"`
}

type MusicRequest struct {
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	ShouldBeMuted bool   `json:"should_be_muted"`
}

type PerformanceRequest struct {
	Title     string `json:"title"`
	Performer string `json:"performer"`
}

// func (m MusicPostRequest) ToDomainMusicPost() entities.MusicPost {
// 	return entities.MusicPost{
// 		Music: entities.Music{
// 			Title:         m.Music.Title,
// 			Artist:        m.Music.Artist,
// 			ShouldBeMuted: m.Music.ShouldBeMuted,
// 		},
// 	}
// }

func (m MusicRequest) ToDomainMusicPost() entities.Music {
	return entities.Music{
		Title:         m.Title,
		Artist:        m.Artist,
		ShouldBeMuted: m.ShouldBeMuted,
	}
}

func (p PerformanceRequest) ToDomainPerformance() entities.Performance {
	return entities.Performance{
		Title:     p.Title,
		Performer: p.Performer,
	}
}

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
