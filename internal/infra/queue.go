package infra

import (
	"github.com/sohosai/ultradonguri-server/internal/core/entity"
	queue_repository "github.com/sohosai/ultradonguri-server/internal/core/repository/queue"
)

func NewQueueRepository() queue_repository.QueueRepository {
	return &queueRepository{}
}

type queueRepository struct {
}

func (r *queueRepository) QueueByID(id string) ([]*entity.Queue, error) {
	var queue []*entity.Queue
	//JSONとかDBから取得する処理を実装する
	// ここではサンプルとして、固定の値を返す
	queue = append(queue, &entity.Queue{
		ID:          id,
		TimetableID: "Sample Timetable",
		Order:       1,
		Name:        "Sample Name",
		QueueType:   "telop",
	})
	return queue, nil
}
