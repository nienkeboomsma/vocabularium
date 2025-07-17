package driving

import (
	"context"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type WorkRepository interface {
	Save(ctx context.Context, db database.Executor, w domain.Work, authorID uuid.UUID) (domain.Work, error)
}
