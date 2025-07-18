package driven

import (
	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/domain"
)

type TextProcessor interface {
	Process([]byte) (*[]domain.WorkWord, *map[uuid.UUID]domain.Word, []string, error)
}
