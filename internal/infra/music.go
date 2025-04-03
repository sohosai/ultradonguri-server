package infra

import (
	"github.com/sohosai/ultradonguri-server/internal/core/entity"
	music_repository "github.com/sohosai/ultradonguri-server/internal/core/repository/music"
)

func NewMusicRepository() music_repository.MusicRepository {
	return &musicRepository{}
}

type musicRepository struct {
}

func (r *musicRepository) MusicByID(id string) ([]*entity.Music, error) {
	var music []*entity.Music
	//JSONとかDBから取得する処理を実装する
	// ここではサンプルとして、固定の値を返す
	music = append(music, &entity.Music{
		ID:            id,
		TimetableID:   "Sample Timetable",
		Order:         1,
		Artist:        "Sample Artist",
		Title:         "Sample Title",
		StreamAllowed: true,
		Note:          "Sample Note",
	})
	return music, nil
}
