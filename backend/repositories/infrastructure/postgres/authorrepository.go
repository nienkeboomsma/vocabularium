package postgres

import (
	"context"
	"database/sql"

	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
)

type AuthorRepository struct {
}

func NewAuthorRepository() *AuthorRepository {
	return &AuthorRepository{}
}

func (wr *AuthorRepository) Save(ctx context.Context, db database.Executor, a domain.Author) (domain.Author, error) {
	q := `
	INSERT INTO author (id, name, modified_at, deleted_at)
	VALUES ($1, $2, DEFAULT, $3)
	ON CONFLICT (name) DO UPDATE
	SET modified_at = DEFAULT, deleted_at = $3
	RETURNING id, name, created_at, modified_at, deleted_at;
	`

	deleted := sql.NullTime{
		Time:  a.Deleted,
		Valid: !a.Deleted.IsZero(),
	}

	var updatedAuthor domain.Author

	err := db.QueryRow(
		ctx,
		q,
		a.ID,
		a.Name,
		deleted,
	).Scan(
		&updatedAuthor.ID,
		&updatedAuthor.Name,
		&updatedAuthor.Created,
		&updatedAuthor.Modified,
		&deleted,
	)
	if err != nil {
		return domain.Author{}, err
	}

	updatedAuthor.Deleted = deleted.Time

	return updatedAuthor, nil
}
