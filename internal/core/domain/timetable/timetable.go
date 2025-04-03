package timetable_domain

import (
	"time"

	music_domain "github.com/sohosai/ultradonguri-server/internal/core/domain/music"
	queue_domain "github.com/sohosai/ultradonguri-server/internal/core/domain/queue"
	"github.com/sohosai/ultradonguri-server/internal/core/entity"
)

type Timetable struct {
	ID        string
	Name      string
	ChannelID string
	StartTime time.Time
	EndTime   time.Time
}

func FromEntity(timetable *entity.TimeTable) *Timetable {
	if timetable == nil {
		return nil
	}

	return &Timetable{
		ID:        timetable.ID,
		Name:      timetable.Name,
		ChannelID: timetable.ChannelID,
		StartTime: timetable.StartTime,
		EndTime:   timetable.EndTime,
	}
}

type TimeTableInfo struct {
	ID        string
	Name      string
	ChannelID string
	StartTime time.Time
	EndTime   time.Time
	Queue     []*queue_domain.Queue
	Music     []*music_domain.Music
}

func FromEntityWithQueueAndMusic(timetable *entity.TimeTable, queues []*queue_domain.Queue, musics []*music_domain.Music) *TimeTableInfo {
	if timetable == nil {
		return nil
	}

	return &TimeTableInfo{
		ID:        timetable.ID,
		Name:      timetable.Name,
		ChannelID: timetable.ChannelID,
		StartTime: timetable.StartTime,
		EndTime:   timetable.EndTime,
		Queue:     queues,
		Music:     musics,
	}
}
