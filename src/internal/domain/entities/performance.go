package entities

type PerformancePost struct {
	Music       Music
	Performance Performance
}

type MusicPost struct {
	Music Music
}

type Music struct {
	Title         string
	Artist        string
	ShouldBeMuted bool
}

type Performance struct {
	Title     string
	Performer string
}
