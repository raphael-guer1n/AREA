package domain

import "time"

type Event struct {
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
}
