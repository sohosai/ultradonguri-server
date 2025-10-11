package responses

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type PerformancesResponse struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Performer   string      `json:"performer"`
	Description string      `json:"description"`
	StartsAt    entities.HM `json:"starts_at"`
	EndsAt      entities.HM `json:"ends_at"`
	Musics      `json:"musics"`
}

type Musics []MusicResponse

type MusicResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	ShouldBeMuted bool   `json:"should_be_muted"`
	Intro         string `json:"intro"`
}

func NewPerformancesResponse(p entities.PerformanceForPerformances) PerformancesResponse {
	return PerformancesResponse{
		ID:          p.ID,
		Title:       p.Title,
		Performer:   p.Performer,
		Description: p.Description,
		StartsAt:    p.StartsAt,
		EndsAt:      p.EndsAt,
		Musics:      NewMusicsResponse(p.Musics),
	}
}

func NewMusicsResponse(musics []entities.MusicForPerformances) Musics {
	res := make(Musics, len(musics))
	for i, m := range musics {
		res[i] = MusicResponse{
			ID:            m.ID,
			Title:         m.Title,
			Artist:        m.Artist,
			ShouldBeMuted: m.ShouldBeMuted,
			Intro:         m.Intro,
		}
	}
	return res
}
