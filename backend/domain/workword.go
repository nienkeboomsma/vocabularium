package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkWord struct {
	ID                        uuid.UUID
	WordID                    uuid.UUID
	WorkID                    uuid.UUID
	WordIndex                 int
	SentenceIndex             int
	OriginalForm              string
	Tag                       string
	MorphoSyntacticalAnalysis string
	Created                   time.Time
	Modified                  time.Time
	Deleted                   time.Time
}
