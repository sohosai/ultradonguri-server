package service

import (
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
)

type TelopService interface {
	SetPerformanceTelop(entities.PerformancePost)
	SetConversionTelop(entities.ConversionPost)
}
