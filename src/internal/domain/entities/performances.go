package entities

type MusicForPerformances struct {
	ID            string
	Title         string
	Artist        string
	ShouldBeMuted bool
	Intro         string
}

type PerformanceForPerformances struct {
	ID          string
	Title       string
	Performer   string
	Description string
	StartsAt    HM
	EndsAt      HM
	Musics      []MusicForPerformances
}
