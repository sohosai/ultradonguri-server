package entities

// type PerformancePost struct {
// 	Music       Music       `json:"music"`
// 	Performance Performance `json:"performance"`
// }

// type Music struct {
// 	Title         string `json:"title"`
// 	Artist        string `json:"artist"`
// 	ShouldBeMuted bool   `json:"should_be_muted"`
// }

// type Performance struct {
// 	Title     string  `json:"title"`
// 	Performer *string `json:"performer"` // null を許すためポインタ
// }

type PerformancePost struct {
	Music       Music
	Performance Performance
}

type Music struct {
	Title         string
	Artist        string
	ShouldBeMuted bool
}

type Performance struct {
	Title     string
	Performer *string // null を許すためポインタ
}
