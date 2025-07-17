package domain

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID       uuid.UUID
	Name     string
	Created  time.Time
	Modified time.Time
	Deleted  time.Time
}
