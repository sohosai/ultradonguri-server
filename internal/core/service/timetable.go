package timetable_service

import (
	"log"

	music_domain "github.com/sohosai/ultradonguri-server/internal/core/domain/music"
	queue_domain "github.com/sohosai/ultradonguri-server/internal/core/domain/queue"
	timetable_domain "github.com/sohosai/ultradonguri-server/internal/core/domain/timetable"
	music_interface "github.com/sohosai/ultradonguri-server/internal/core/repository/music"
	queue_interface "github.com/sohosai/ultradonguri-server/internal/core/repository/queue"
	queue_repository "github.com/sohosai/ultradonguri-server/internal/core/repository/queue"
	timetable_interface "github.com/sohosai/ultradonguri-server/internal/core/repository/timetable"
)

type TimetableService struct {
	queueRepository     queue_interface.QueueRepository
	timetableRepository timetable_interface.TimetableRepository
	musicRepository     music_interface.MusicRepository
}

func New(
	queueRepository queue_interface.QueueRepository,
	timetableRepository timetable_interface.TimetableRepository,
	musicRepository music_interface.MusicRepository,
) *TimetableService {
	return &TimetableService{
		queueRepository:     queueRepository,
		timetableRepository: timetableRepository,
		musicRepository:     musicRepository,
	}
}
func NewQueueService(queueRepository queue_repository.QueueRepository) *TimetableService {
	return &TimetableService{
		queueRepository: queueRepository,
	}
}

func (s *TimetableService) TimetableByID(id string) (*timetable_domain.TimeTableInfo, error) {
	timetableEntity, err := s.timetableRepository.TimetableByID(id)
	log.Println("timetableEntity", timetableEntity)
	if err != nil {
		return nil, err
	}
	if timetableEntity == nil {
		return nil, nil
	}
	timetable := timetable_domain.FromEntity(timetableEntity)
	log.Println("timetable", timetable)
	if timetable == nil {
		return nil, nil
	}

	queueEntity, err := s.queueRepository.QueueByID(timetable.ID)
	log.Println("queueEntity", queueEntity)
	if err != nil {
		return nil, err
	}
	if queueEntity == nil {
		return nil, nil
	}

	var queue []*queue_domain.Queue
	for _, q := range queueEntity {
		queue = append(queue, queue_domain.FromEntity(q))
	}
	if queue == nil {
		return nil, nil
	}

	musicEntity, err := s.musicRepository.MusicByID(timetable.ID)
	if err != nil {
		return nil, err
	}
	if musicEntity == nil {
		return nil, nil
	}
	var music []*music_domain.Music
	for _, m := range musicEntity {
		music = append(music, music_domain.FromEntity(m))
	}
	if music == nil {
		return nil, nil
	}

	timetableInfo := timetable_domain.FromEntityWithQueueAndMusic(timetableEntity, queue, music)
	if timetableInfo == nil {
		return nil, nil
	}
	return timetableInfo, nil

}
