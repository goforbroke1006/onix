package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewSourceRepository creates data exchange object with db.
func NewSourceRepository(conn *pgxpool.Pool) *sourceRepository { // nolint:revive,golint
	return &sourceRepository{
		conn: conn,
	}
}

var _ domain.SourceRepository = &sourceRepository{} // nolint:exhaustivestruct

type sourceRepository struct {
	conn *pgxpool.Pool
}

func (repo sourceRepository) Create(title string, kind domain.SourceType, address string) (int64, error) {
	query := fmt.Sprintf(
		"INSERT INTO source (title, kind, address) VALUES ('%s', '%s', '%s') RETURNING id;",
		title, kind, address,
	)

	var (
		rows pgx.Rows
		err  error
	)

	if rows, err = repo.conn.Query(context.TODO(), query); err != nil {
		return 0, errors.Wrap(err, "can't exec query")
	}
	defer rows.Close()

	var identifier int64
	if rows.Next() {
		if err := rows.Scan(&identifier); err != nil {
			return 0, errors.Wrap(err, "can't scan source row")
		}

		return identifier, nil
	}

	return 0, domain.ErrNotFound
}

func (repo sourceRepository) Get(identifier int64) (*domain.Source, error) {
	query := fmt.Sprintf(
		`
SELECT title, kind, address 
FROM source 
WHERE id = %d
;`,
		identifier,
	)

	var (
		rows pgx.Rows
		err  error
	)

	if rows, err = repo.conn.Query(context.TODO(), query); err != nil {
		return nil, errors.Wrap(err, "can't exec query")
	}
	defer rows.Close()

	var (
		title   string
		kind    domain.SourceType
		address string
	)

	if rows.Next() {
		if err := rows.Scan(&title, &kind, &address); err != nil {
			return nil, errors.Wrap(err, "can't scan source row")
		}

		release := domain.Source{
			ID:      identifier,
			Title:   title,
			Type:    kind,
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

	var (
		rows pgx.Rows
		err  error
	)

	if rows, err = repo.conn.Query(context.TODO(), query); err != nil {
		return nil, errors.Wrap(err, "can't exec query")
	}
	defer rows.Close()

	var (
		identifier int64
		title      string
		kind       domain.SourceType
		address    string
	)

	result := make([]domain.Source, 0, len(rows.RawValues()))

	for rows.Next() {
		if err := rows.Scan(&identifier, &title, &kind, &address); err != nil {
			return nil, errors.Wrap(err, "can't scan source row")
		}

		result = append(result, domain.Source{
			ID:      identifier,
			Title:   title,
			Type:    kind,
			Address: address,
		})
	}

	return result, nil
}
