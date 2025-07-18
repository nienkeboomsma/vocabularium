package driving

import (
	"context"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/domain"
)

type WorkPersister interface {
	Persist(ctx context.Context, author domain.Author, work domain.Work, words *map[uuid.UUID]domain.Word, workWords *[]domain.WorkWord) error
}
