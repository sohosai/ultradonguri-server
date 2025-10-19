package entities

type AppError struct {
	Message string
	Kind    int
}

const (
	NoConnectionToOBS = iota
	NoConnectionToViewer
	InvalidPerformancesJson
	InvalidFormat
	CannotConversion
	CannotForceMute
	CannotChangeState
)
