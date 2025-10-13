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

// type MusicForPerformances struct {
// 	ID            string `json:"id"`
// 	Title         string `json:"title"`
// 	Artist        string `json:"artist"`
// 	ShouldBeMuted bool   `json:"should_be_muted"`
// 	Intro         string `json:"intro"`
// }

// type PerformanceForPerformances struct {
// 	ID          string                 `json:"id"`
// 	Title       string                 `json:"title"`
// 	Performer   string                 `json:"performer"`
// 	Description string                 `json:"description"`
// 	StartsAt    HM                     `json:"starts_at"`
// 	EndsAt      HM                     `json:"ends_at"`
// 	Musics      []MusicForPerformances `json:"musics"`
// }
