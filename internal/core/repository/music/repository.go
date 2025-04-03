package music_interface

import "github.com/sohosai/ultradonguri-server/internal/core/entity"

type MusicRepository interface {
	MusicByID(id string) ([]*entity.Music, error)
}
