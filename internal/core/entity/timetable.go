package entity

import (
	"time"
)

type TimeTable struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ChannelID string    `json:"channel_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
