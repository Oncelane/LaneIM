package model

import (
	"time"
)

// Person represents a row in person table.
// Field names are converted to snake case by default, no need to add special tags.
// A field will not be persisted by adding the `db:"-"` tag or making it unexported.
type ChatMessage struct {
	GroupID   int64
	UserID    int64
	UserSeq   int64
	Timestamp time.Time
	MessageID int64
	Content   []byte
}
