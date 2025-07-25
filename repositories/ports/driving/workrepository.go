package driving

import (
	"context"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
)

type WorkRepository interface {
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context) ([]domain.Work, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Work, error)
	Save(ctx context.Context, db database.Executor, w domain.Work, authorID uuid.UUID) (domain.Work, error)
}
