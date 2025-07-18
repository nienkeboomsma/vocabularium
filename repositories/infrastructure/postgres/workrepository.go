package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
)

type WorkRepository struct {
	db *database.Client
}

func NewWorkRepository(db *database.Client) *WorkRepository {
	return &WorkRepository{db: db}
}

func (wr *WorkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := wr.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	q := `
	UPDATE work
	SET deleted_at = NOW()
	WHERE id = $1
	AND deleted_at IS NULL;
	`

	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to execute work query: %w", err)
	}

	q = `
	UPDATE work_word
	SET deleted_at = NOW()
	WHERE work_id = $1
	AND deleted_at IS NULL;
	`

	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to execute work_word query: %w", err)
	}

	q = `
	UPDATE author
	SET deleted_at = NOW()
	WHERE id = (
	    SELECT author_id FROM work WHERE id = $1
	)
	AND NOT EXISTS (
	    SELECT 1 FROM work
	    WHERE author_id = (
	        SELECT author_id FROM work WHERE id = $1
	    )
	    AND deleted_at IS NULL
	)
	AND deleted_at IS NULL;
	`

	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to execute author query: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (wr *WorkRepository) Get(ctx context.Context) ([]domain.Work, error) {
	q := `
	SELECT w.id, a.id, a.name, w.title
	FROM work w
	JOIN author a
	ON a.id = w.author_id
	WHERE w.deleted_at IS NULL;
	`

	works := []domain.Work{}

	rows, err := wr.db.Pool.Query(ctx, q)
	if err != nil {
		return []domain.Work{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		work := domain.Work{}

		err = rows.Scan(&work.ID, &work.Author.ID, &work.Author.Name, &work.Title)
		if err != nil {
			return []domain.Work{}, fmt.Errorf("failed to scan row: %w", err)
		}

		works = append(works, work)
	}

	err = rows.Err()
	if err != nil {
		return []domain.Work{}, fmt.Errorf("failed to read rows: %w", err)
	}

	return works, nil
}

func (wr *WorkRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Work, error) {
	q := `
	SELECT w.id, a.name, w.title, w.created_at, w.modified_at, w.deleted_at
	FROM work w
	JOIN author a
	ON a.id = w.author_id
	WHERE w.id = $1
	AND w.deleted_at IS NULL;
	`

	var work domain.Work
	var deleted sql.NullTime

	err := wr.db.Pool.QueryRow(
		ctx,
		q,
		id,
	).Scan(
		&work.ID,
		&work.Author.Name,
		&work.Title,
		&work.Created,
		&work.Modified,
		&deleted,
	)
	if err != nil {
		return domain.Work{}, err
	}

	work.Deleted = deleted.Time

	return work, nil
}

func (wr *WorkRepository) Save(ctx context.Context, db database.Executor, w domain.Work, authorID uuid.UUID) (domain.Work, error) {
	q := `
	INSERT INTO work (id, author_id, title, modified_at, deleted_at)
	VALUES ($1, $2, $3, DEFAULT, $4)
	ON CONFLICT (author_id, title) DO UPDATE
	SET modified_at = DEFAULT, deleted_at = $4
	RETURNING id, author_id, title, created_at, modified_at, deleted_at;
	`

	deleted := sql.NullTime{
		Time:  w.Deleted,
		Valid: !w.Deleted.IsZero(),
	}

	var updatedWork domain.Work

	err := db.QueryRow(
		ctx,
		q,
		w.ID,
		authorID,
		w.Title,
		deleted,
	).Scan(
		&updatedWork.ID,
		&authorID,
		&updatedWork.Title,
		&updatedWork.Created,
		&updatedWork.Modified,
		&deleted,
	)
	if err != nil {
		return domain.Work{}, err
	}

	updatedWork.Deleted = deleted.Time

	return updatedWork, nil
}
