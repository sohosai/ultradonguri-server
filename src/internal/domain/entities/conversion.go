package entities

type ConversionPost struct {
	NextPerformances []NextPerformance
}

type NextPerformance struct {
	Title       string
	Performer   string
	Description string
}
