package entities

type MusicForPerformances struct {
	ID     string
	Title  string
	Artist string
	Intro  string
}

type PerformanceForPerformances struct {
	ID          string
	Title       string
	Performer   string
	Description string
	StartsAt    ISOTime
	EndsAt      ISOTime
	Musics      []MusicForPerformances
}
