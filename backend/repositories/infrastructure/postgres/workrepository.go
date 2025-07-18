package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type WorkRepository struct {
	db *database.Client
}

func NewWorkRepository(db *database.Client) *WorkRepository {
	return &WorkRepository{db: db}
}

func (wr *WorkRepository) Get(ctx context.Context) ([]domain.Work, error) {
	q := `
	SELECT w.id, a.id, a.name, w.title, w.type
	FROM work w
	JOIN author a
	ON a.id = w.author_id;
	`

	works := []domain.Work{}

	rows, err := wr.db.Pool.Query(ctx, q)
	if err != nil {
		return []domain.Work{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		work := domain.Work{}

		err = rows.Scan(&work.ID, &work.Author.ID, &work.Author.Name, &work.Title, &work.Type)
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
	SELECT w.id, a.name, w.title, w.type, w.created_at, w.modified_at, w.deleted_at
	FROM work w
	JOIN author a
	ON a.id = w.author_id
	WHERE w.id = $1;
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
		&work.Type,
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
	INSERT INTO work (id, author_id, title, type, modified_at, deleted_at)
	VALUES ($1, $2, $3, $4, DEFAULT, $5)
	ON CONFLICT (author_id, title) DO UPDATE
	SET type = $4, modified_at = DEFAULT, deleted_at = $5
	RETURNING id, author_id, title, type, created_at, modified_at, deleted_at;
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
		w.Type,
		deleted,
	).Scan(
		&updatedWork.ID,
		&authorID,
		&updatedWork.Title,
		&updatedWork.Type,
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
