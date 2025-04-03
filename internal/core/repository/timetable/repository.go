package timetable_interface

import (
	"github.com/sohosai/ultradonguri-server/internal/core/entity"
)

type TimetableRepository interface {
	TimetableByID(id string) (*entity.TimeTable, error)
}
