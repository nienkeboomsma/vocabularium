package postgres

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
	"github.com/nienkeboomsma/vocabularium/repositories/ports/driving"
)

type WorkPersister struct {
	db                 *database.Client
	authorRepository   driving.AuthorRepository
	workRepository     driving.WorkRepository
	wordRepository     driving.WordRepository
	workWordRepository driving.WorkWordRepository
}

func NewWorkPersister(
	db *database.Client,
	authorRepository driving.AuthorRepository,
	workRepository driving.WorkRepository,
	wordRepository driving.WordRepository,
	workWordRepository driving.WorkWordRepository,
) *WorkPersister {
	return &WorkPersister{
		db:                 db,
		authorRepository:   authorRepository,
		workRepository:     workRepository,
		wordRepository:     wordRepository,
		workWordRepository: workWordRepository,
	}
}

func (wp *WorkPersister) Persist(ctx context.Context, author domain.Author, work domain.Work, words *map[uuid.UUID]domain.Word, workWords *[]domain.WorkWord) error {
	tx, err := wp.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %w", err)
	}

	updatedAuthor, err := wp.authorRepository.Save(ctx, tx, author)

	updatedWork, err := wp.workRepository.Save(ctx, tx, work, updatedAuthor.ID)

	for _, word := range *words {
		_, err := wp.wordRepository.Insert(ctx, tx, word)
		if err != nil {
			return fmt.Errorf("failed to save word: %w", err)
		}
	}

	for _, workWord := range *workWords {
		if reflect.DeepEqual(workWord.WordID, uuid.UUID{}) {
			continue
		}

		_, err := wp.workWordRepository.Save(ctx, tx, workWord, updatedWork.ID)
		if err != nil {
			return fmt.Errorf("failed to save work_word: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}

	return nil
}
