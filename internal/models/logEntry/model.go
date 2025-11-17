package logentry

import (
	"time"

	"github.com/julian-richter/ApiTemplate/internal/models"
)

// LogEntry represents an application log entry.
type LogEntry struct {
	ID        int       `json:"id" db:"id"`
	Level     string    `json:"level" db:"level"`
	Message   string    `json:"message" db:"message"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func (l *LogEntry) GetID() int64 {
	return int64(l.ID)
}

func (l *LogEntry) SetID(id int64) {
	l.ID = int(id)
}

var _ models.Entity = (*LogEntry)(nil)
