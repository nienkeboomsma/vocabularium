package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type WordRepository struct {
	db *database.Client
}

func NewWordRepository(db *database.Client) *WordRepository {
	return &WordRepository{db: db}
}

func (wr *WordRepository) GetFrequencyListByAuthorID(ctx context.Context, authorID uuid.UUID) (*[]domain.WordInWork, error) {
	q := `
	SELECT w.id, w.lemma_rich, w.translation, w.known, COUNT(ww.word_id)
	FROM work_word ww
	JOIN word w
	ON w.id = ww.word_id
	JOIN work
	ON work.id = ww.work_id
	JOIN author a
	ON a.id = work.author_id
	WHERE a.id = $1
	GROUP BY w.id, w.lemma_rich, w.known
	ORDER BY COUNT(ww.word_id) DESC;
	`

	return wr.getWordListByWorkID(ctx, q, authorID)
}

func (wr *WordRepository) GetFrequencyListByWorkID(ctx context.Context, workID uuid.UUID) (*[]domain.WordInWork, error) {
	q := `
	SELECT w.id, w.lemma_rich, w.translation, w.known, COUNT(ww.word_id)
	FROM work_word ww
	JOIN word w
	ON w.id = ww.word_id
	WHERE ww.work_id = $1
	GROUP BY w.id, w.lemma_rich, w.known
	ORDER BY COUNT(ww.word_id) DESC;
	`

	return wr.getWordListByWorkID(ctx, q, workID)
}

func (wr *WordRepository) GetGlossaryByWorkID(ctx context.Context, workID uuid.UUID) (*[]domain.WordInWork, error) {
	q := `
	SELECT w.id, w.lemma_rich, w.translation, w.known, COUNT(ww.word_id) OVER (PARTITION BY w.id) AS word_count
	FROM work_word ww
	JOIN word w
	ON w.id = ww.word_id
	WHERE ww.work_id = $1
	ORDER BY ww.word_index ASC;
	`

	return wr.getWordListByWorkID(ctx, q, workID)
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

func (wr *WordRepository) ToggleKnownStatus(ctx context.Context, wordID uuid.UUID) (domain.Word, error) {
	q := `
	UPDATE word
	SET known = NOT known
	WHERE id = $1
	RETURNING *;
	`

	var word domain.Word

	err := wr.db.Pool.QueryRow(ctx, q, wordID).Scan(
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

func (wr *WordRepository) getWordListByWorkID(ctx context.Context, q string, args ...any) (*[]domain.WordInWork, error) {
	words := []domain.WordInWork{}

	rows, err := wr.db.Pool.Query(ctx, q, args...)
	if err != nil {
		return &[]domain.WordInWork{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		word := domain.WordInWork{}

		err = rows.Scan(&word.ID, &word.LemmaRich, &word.Translation, &word.Known, &word.Count)
		if err != nil {
			return &[]domain.WordInWork{}, fmt.Errorf("failed to scan row: %w", err)
		}

		words = append(words, word)
	}

	err = rows.Err()
	if err != nil {
		return &[]domain.WordInWork{}, fmt.Errorf("failed to read rows: %w", err)
	}

	return &words, nil
}
