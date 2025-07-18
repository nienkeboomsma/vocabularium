package driving

import (
	"context"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
)

type WorkWordRepository interface {
	Save(ctx context.Context, db database.Executor, ww domain.WorkWord, workID uuid.UUID) (domain.WorkWord, error)
}
