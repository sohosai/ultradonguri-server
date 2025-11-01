package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type PerformancePostRequest struct {
	Music       MusicRequest       `json:"music" binding:"required"`
	Performance PerformanceRequest `json:"performance" binding:"required"`
}

type MusicPostRequest struct {
	Music MusicRequest `json:"music" binding:"required"`
}

type OnlyPerformancePostRequest struct {
	Performance PerformanceRequest `json:"performance" binding:"required"`
}

type MusicRequest struct {
	Title  *string `json:"title" binding:"required"`
	Artist *string `json:"artist" binding:"required"`
	// binding:"required"の影響でfalseを渡したときにエラーが出るので、bool型ではなくポインタを渡す
	ShouldBeMuted *bool `json:"should_be_muted" binding:"required"`
}

type PerformanceRequest struct {
	Title     string `json:"title" binding:"required"`
	Performer string `json:"performer" binding:"required"`
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
		Title:         *m.Title,
		Artist:        *m.Artist,
		ShouldBeMuted: *m.ShouldBeMuted,
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
			Title:         *pp.Music.Title,
			Artist:        *pp.Music.Artist,
			ShouldBeMuted: *pp.Music.ShouldBeMuted,
		},
		Performance: entities.Performance{
			Title:     pp.Performance.Title,
			Performer: pp.Performance.Performer,
		},
	}
}
