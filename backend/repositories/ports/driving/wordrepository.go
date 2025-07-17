package driving

import (
	"context"

	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type WordRepository interface {
	Save(ctx context.Context, db database.Executor, w domain.Word) (domain.Word, error)
}
