package spec

type ConversionPost struct {
	NextPerformances []NextPerformance `json:"next_performances"`
}

type NextPerformance struct {
	Title       string `json:"title"`
	Performer   string `json:"performer"`
	Description string `json:"description"`
}
