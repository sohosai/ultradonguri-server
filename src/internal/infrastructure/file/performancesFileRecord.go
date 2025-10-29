package file

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type PerformancesRecords []PerformancesRecord

type PerformancesRecord struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Performer   string           `json:"performer"`
	Description string           `json:"description"`
	StartsAt    entities.ISOTime `json:"starts_at"`
	EndsAt      entities.ISOTime `json:"ends_at"`
	Musics      `json:"musics"`
}

type Musics []MusicRecord

type MusicRecord struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	ShouldBeMuted bool   `json:"should_be_muted"`
	Intro         string `json:"intro"`
}

func (perf PerformancesRecords) ToDomainPerformanceForPerformances() []entities.PerformanceForPerformances {
	res := make([]entities.PerformanceForPerformances, len(perf))
	for i, p := range perf {
		res[i] = entities.PerformanceForPerformances{
			ID:          p.ID,
			Title:       p.Title,
			Performer:   p.Performer,
			Description: p.Description,
			StartsAt:    p.StartsAt,
			EndsAt:      p.EndsAt,
			Musics:      p.Musics.ToDomainMusicsForPerformances(),
		}
	}
	return res
}

func (m Musics) ToDomainMusicsForPerformances() []entities.MusicForPerformances {
	res := make([]entities.MusicForPerformances, len(m))
	for i, music := range m {
		res[i] = entities.MusicForPerformances{
			ID:            music.ID,
			Title:         music.Title,
			Artist:        music.Artist,
			ShouldBeMuted: music.ShouldBeMuted,
			Intro:         music.Intro,
		}
	}
	return res
}
