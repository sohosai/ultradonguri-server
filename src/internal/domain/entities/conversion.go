package entities

type ConversionPost struct {
	NextPerformances []NextPerformance
}

type NextPerformance struct {
	Title       string
	Performer   string
	Description string
	StartsAt    HM
}

type CMState struct {
	IsCMMode bool
}
