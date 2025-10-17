package entities

type TelopType string

const (
	TelopTypePerformance TelopType = "performance"
	TelopTypeConversion  TelopType = "conversion"
	TelopTypeEmpty       TelopType = "empty"
)

type TelopMessage struct {
	Type            TelopType
	PerformanceData *PerformancePost
	ConversionData  *ConversionPost
}
