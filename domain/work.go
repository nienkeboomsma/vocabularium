package domain

import (
	"time"

	"github.com/google/uuid"
)

type Work struct {
	ID       uuid.UUID
	Title    string
	Author   Author
	Created  time.Time
	Modified time.Time
	Deleted  time.Time
}
