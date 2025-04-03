package queue_domain

import (
	"slices"

	"github.com/sohosai/ultradonguri-server/internal/core/entity"
	"github.com/sohosai/ultradonguri-server/internal/libs/enum"
)

type Queue struct {
	ID          string
	TimetableID string
	Name        string
	Order       int
	QueueType   enum.EnumyTypes
}

type QueueType string

const (
	QueueTypeTelop QueueType = "telop"
	QueueTypeOther QueueType = "other"
)

var validQueueType = []QueueType{
	QueueTypeTelop,
	QueueTypeOther,
}

func (q QueueType) Valid() bool {
	return slices.Contains(validQueueType, q)
}

func (q QueueType) Get() string {
	return string(q)
}

func FromEntity(queue *entity.Queue) *Queue {
	if queue == nil {
		return nil
	}
	queueType := QueueType(queue.QueueType)
	if !queueType.Valid() {
		queueType = QueueTypeOther
	}

	return &Queue{
		ID:        queue.ID,
		Name:      queue.Name,
		QueueType: queueType,
	}
}
