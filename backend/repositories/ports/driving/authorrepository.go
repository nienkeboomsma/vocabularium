package driving

import (
	"context"

	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type AuthorRepository interface {
	Save(ctx context.Context, db database.Executor, a domain.Author) (domain.Author, error)
}
