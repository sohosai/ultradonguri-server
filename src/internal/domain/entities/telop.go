package entities

type TelopType string

const (
	TelopTypePerformance TelopType = "performance"
	TelopTypeConversion  TelopType = "conversion"
	TelopTypeEmpty       TelopType = "empty"
)

type TelopMessage struct {
	Type            TelopType        `json:"type"`
	PerformanceData *PerformancePost `json:"performance_data,omitempty"`
	ConversionData  *ConversionPost  `json:"conversion_data,omitempty"`
}
