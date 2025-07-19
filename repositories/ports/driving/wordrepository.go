package driving

import (
	"context"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
)

type WordRepository interface {
	GetFrequencyList(ctx context.Context) (*[]domain.WordInWork, error)
	GetFrequencyListByAuthorID(ctx context.Context, authorID uuid.UUID) (*[]domain.WordInWork, error)
	GetFrequencyListByWorkID(ctx context.Context, workID uuid.UUID) (*[]domain.WordInWork, error)
	GetGlossaryByWorkID(ctx context.Context, workID uuid.UUID) (*[]domain.WordInWork, error)
	Insert(ctx context.Context, db database.Executor, w domain.Word) (domain.Word, error)
	ToggleKnownStatus(ctx context.Context, wordID uuid.UUID) (domain.Word, error)
}
