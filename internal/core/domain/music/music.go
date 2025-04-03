package music_domain

import "github.com/sohosai/ultradonguri-server/internal/core/entity"

type Music struct {
	ID            string
	TimetableID   string
	Order         int
	Artist        string
	Title         string
	StreamAllowed bool
	Note          string
}

func FromEntity(music *entity.Music) *Music {
	if music == nil {
		return nil
	}

	return &Music{
		ID:            music.ID,
		TimetableID:   music.TimetableID,
		Order:         music.Order,
		Artist:        music.Artist,
		Title:         music.Title,
		StreamAllowed: music.StreamAllowed,
		Note:          music.Note,
	}
}
