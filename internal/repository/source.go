package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/goforbroke1006/onix/domain"
)

// NewSourceRepository creates data exchange object with db
func NewSourceRepository(conn *pgxpool.Pool) *sourceRepository { // nolint:golint
	return &sourceRepository{
		conn: conn,
	}
}

var (
	_ domain.SourceRepository = &sourceRepository{}
)

type sourceRepository struct {
	conn *pgxpool.Pool
}

func (repo sourceRepository) Create(title string, kind domain.SourceType, address string) (int64, error) {
	query := fmt.Sprintf(
		"INSERT INTO source (title, kind, address) VALUES ('%s', '%s', '%s') RETURNING id;",
		title, kind, address,
	)
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id int64
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}

		return id, nil
	}

	return 0, domain.ErrNotFound
}

func (repo sourceRepository) Get(id int64) (*domain.Source, error) {
	query := fmt.Sprintf(
		`
SELECT title, kind, address 
FROM source 
WHERE id = %d
;`,
		id,
	)
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		title   string
		kind    domain.SourceType
		address string
	)

	if rows.Next() {
		if err := rows.Scan(&title, &kind, &address); err != nil {
			return nil, err
		}

		release := domain.Source{
			ID:      id,
			Title:   title,
			Kind:    kind,
			Address: address,
		}

		return &release, nil
	}

	return nil, domain.ErrNotFound
}

func (repo sourceRepository) GetAll() ([]domain.Source, error) {
	query := `
		SELECT id, title, kind, address
		FROM source 
		ORDER BY id ASC
	;`
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id      int64
		title   string
		kind    domain.SourceType
		address string
	)

	result := make([]domain.Source, 0, len(rows.RawValues()))
	for rows.Next() {
		if err := rows.Scan(&id, &title, &kind, &address); err != nil {
			return nil, err
		}
		result = append(result, domain.Source{
			ID:      id,
			Title:   title,
			Kind:    kind,
			Address: address,
		})
	}

	return result, nil
}
