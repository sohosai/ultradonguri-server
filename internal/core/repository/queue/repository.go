package queue_interface

import (
	"github.com/sohosai/ultradonguri-server/internal/core/entity"
)

type QueueRepository interface {
	QueueByID(id string) ([]*entity.Queue, error)
}
