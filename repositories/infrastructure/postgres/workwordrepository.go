package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
)

type WorkWordRepository struct {
}

func NewWorkWordRepository() *WorkWordRepository {
	return &WorkWordRepository{}
}

func (wr *WorkWordRepository) Save(ctx context.Context, db database.Executor, ww domain.WorkWord, workID uuid.UUID) (domain.WorkWord, error) {
	q := `
	INSERT INTO work_word (id, work_id, word_id, word_index, sentence_index, original_form, tag, morph_analysis, modified_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, DEFAULT, $9)
	ON CONFLICT (work_id, word_index) DO UPDATE
	SET sentence_index = $5, original_form = $6, tag = $7, morph_analysis = $8, modified_at = DEFAULT, deleted_at = $9
	RETURNING id, work_id, word_id, word_index, sentence_index, original_form, tag, morph_analysis, created_at, modified_at, deleted_at;
	`

	created := sql.NullTime{}
	modified := sql.NullTime{}
	deleted := sql.NullTime{
		Time:  ww.Deleted,
		Valid: !ww.Deleted.IsZero(),
	}

	var updatedWorkWord domain.WorkWord

	err := db.QueryRow(
		ctx,
		q,
		uuid.New(),
		workID,
		ww.WordID,
		ww.WordIndex,
		ww.SentenceIndex,
		ww.OriginalForm,
		ww.Tag,
		ww.MorphoSyntacticalAnalysis,
		deleted,
	).Scan(
		&updatedWorkWord.ID,
		&workID,
		&updatedWorkWord.WordID,
		&updatedWorkWord.WordIndex,
		&updatedWorkWord.SentenceIndex,
		&updatedWorkWord.OriginalForm,
		&updatedWorkWord.Tag,
		&updatedWorkWord.MorphoSyntacticalAnalysis,
		&created,
		&modified,
		&deleted,
	)
	if err != nil {
		return domain.WorkWord{}, err
	}

	updatedWorkWord.Created = created.Time
	updatedWorkWord.Modified = modified.Time
	updatedWorkWord.Deleted = deleted.Time

	return updatedWorkWord, nil
}
