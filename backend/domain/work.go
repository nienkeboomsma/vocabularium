package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkType string

const (
	WorkTypeProse WorkType = "prose"
	WorkTypeVerse WorkType = "verse"
)

type Work struct {
	ID       uuid.UUID
	Title    string
	Author   Author
	Type     WorkType
	Created  time.Time
	Modified time.Time
	Deleted  time.Time
}
