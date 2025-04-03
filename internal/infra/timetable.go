package infra

import (
	"github.com/sohosai/ultradonguri-server/internal/core/entity"
	timetable_interface "github.com/sohosai/ultradonguri-server/internal/core/repository/timetable"
)

func NewTimetableRepository() timetable_interface.TimetableRepository {
	return &timetableRepository{}
}

type timetableRepository struct {
}

func (r *timetableRepository) TimetableByID(id string) (*entity.TimeTable, error) {
	var timetable *entity.TimeTable
	//JSONとかDBから取得する処理を実装する
	// ここではサンプルとして、固定の値を返す
	timetable = &entity.TimeTable{
		ID:   id,
		Name: "Sample Timetable",
	}
	return timetable, nil
}
