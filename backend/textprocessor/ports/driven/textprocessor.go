package driven

import "github.com/nienkeboomsma/collatinus/domain"

type TextProcessor interface {
	Process([]byte) (*[]domain.Word, error)
}
