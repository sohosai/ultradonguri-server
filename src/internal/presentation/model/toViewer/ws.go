package toViewer

import "github.com/sohosai/ultradonguri-server/internal/domain/entities"

//websocketで送るjsonについて

type WsTelopMessage struct {
	Type            entities.TelopType `json:"type"`
	PerformanceData *WsPerformancePost `json:"performance_data,omitempty"`
	ConversionData  *WsConversionPost  `json:"conversion_data,omitempty"`
}

type WsPerformancePost struct {
	Music       WsMusic       `json:"music"`
	Performance WsPerformance `json:"performance"`
}

type WsMusic struct {
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	ShouldBeMuted bool   `json:"should_be_muted"`
}

type WsPerformance struct {
	Title     string  `json:"title"`
	Performer *string `json:"performer"`
}

type WsConversionPost struct {
	NextPerformances []WsNextPerformance `json:"next_performances"`
}

type WsNextPerformance struct {
	Title       string `json:"title"`
	Performer   string `json:"performer"`
	Description string `json:"description"`
}

func NewWsTelopMessage(t entities.TelopType, p *entities.PerformancePost, c *entities.ConversionPost) WsTelopMessage {
	var wp *WsPerformancePost
	var wc *WsConversionPost

	if p != nil {
		wp = &WsPerformancePost{
			Music: WsMusic{
				Title:         p.Music.Title,
				Artist:        p.Music.Artist,
				ShouldBeMuted: p.Music.ShouldBeMuted,
			},
			Performance: WsPerformance{
				Title:     p.Performance.Title,
				Performer: p.Performance.Performer,
			},
		}
	}

	if c != nil {
		wc = &WsConversionPost{
			NextPerformances: make([]WsNextPerformance, len(c.NextPerformances)),
		}
		for i, np := range c.NextPerformances {
			wc.NextPerformances[i] = WsNextPerformance{
				Title:       np.Title,
				Performer:   np.Performer,
				Description: np.Description,
			}
		}
	}

	return WsTelopMessage{
		Type:            t,
		PerformanceData: wp,
		ConversionData:  wc,
	}
}
