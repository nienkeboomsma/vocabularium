package driving

import (
	"context"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type WorkWordRepository interface {
	Save(ctx context.Context, db database.Executor, ww domain.WorkWord, workID uuid.UUID) (domain.WorkWord, error)
}
