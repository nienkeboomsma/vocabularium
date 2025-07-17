package postgres

import (
	"context"

	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type WordRepository struct {
}

func NewWordRepository() *WordRepository {
	return &WordRepository{}
}

func (wr *WordRepository) Save(ctx context.Context, db database.Executor, w domain.Word) (domain.Word, error) {
	q := `
	INSERT INTO word (id, lemma_raw, lemma_rich, translation, lasla_frequency, known, modified_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, DEFAULT, $7)
	ON CONFLICT (lemma_raw, lemma_rich) DO UPDATE
	SET translation = $4, known = $6, modified_at = DEFAULT, deleted_at = $7
	RETURNING id, lemma_raw, lemma_rich, translation, lasla_frequency, known, created_at, modified_at, deleted_at;
	`

	var word domain.Word

	err := db.QueryRow(
		ctx,
		q,
		w.ID,
		w.LemmaRaw,
		w.LemmaRich,
		w.Translation,
		w.FrequencyInLASLA,
		w.Known,
		w.Deleted,
	).Scan(
		&word.ID,
		&word.LemmaRaw,
		&word.LemmaRich,
		&word.Translation,
		&word.FrequencyInLASLA,
		&word.Known,
		&word.Created,
		&word.Modified,
		&word.Deleted,
	)
	if err != nil {
		return domain.Word{}, err
	}

	return word, nil
}
