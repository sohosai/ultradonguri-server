package requests

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

type DisplayCopyrightRequest struct {
	IsDisplay *bool `json:"is_displayed_copyright" binding:"required"`
}

func (d DisplayCopyrightRequest) ToDomainCopyright() entities.DisplayCopyright {
	return entities.DisplayCopyright{
		IsDisplay: *d.IsDisplay,
	}
}
