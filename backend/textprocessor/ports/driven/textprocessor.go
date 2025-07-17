package driven

import (
	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/domain"
)

type TextProcessor interface {
	Process([]byte) (*[]domain.WorkWord, *map[uuid.UUID]domain.Word, error)
}
