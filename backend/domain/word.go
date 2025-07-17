package domain

import (
	"time"

	"github.com/google/uuid"
)

type Word struct {
	ID               uuid.UUID
	LemmaRaw         string
	LemmaRich        string
	Translation      string
	FrequencyInLASLA int
	Known            bool
	Created          time.Time
	Modified         time.Time
	Deleted          time.Time
}
