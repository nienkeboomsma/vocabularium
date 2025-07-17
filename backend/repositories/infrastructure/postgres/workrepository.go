package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type WorkRepository struct {
}

func NewWorkRepository() *WorkRepository {
	return &WorkRepository{}
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
