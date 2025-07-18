package driving

import (
	"context"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
)

type AuthorRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Author, error)
	Save(ctx context.Context, db database.Executor, a domain.Author) (domain.Author, error)
}
