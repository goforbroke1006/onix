package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewSourceRepository creates data exchange object with db.
func NewSourceRepository(conn *pgxpool.Pool) domain.SourceRepository { // nolint:golint
	return &sourceRepository{conn: conn}
}

var _ domain.SourceRepository = (*sourceRepository)(nil)

type sourceRepository struct {
	conn *pgxpool.Pool
}

func (repo sourceRepository) Create(
	ctx context.Context,
	id string,
	kind domain.SourceType,
	address string,
) error {
	const query = `
		INSERT INTO source (id, kind, address) 
		VALUES (:id, :kind, :address);`

	if _, execErr := repo.conn.Exec(ctx, query, id, kind, address); execErr != nil {
		return errors.Wrap(execErr, "can't exec query")
	}

	return nil
}

func (repo sourceRepository) Get(ctx context.Context, id string) (domain.Source, error) {
	const query = `
		SELECT kind, address 
		FROM source 
		WHERE id = :id;
		`
	rows, rowsErr := repo.conn.Query(ctx, query, id)
	if rowsErr != nil {
		return domain.Source{}, errors.Wrap(rowsErr, "can't exec query")
	}
	defer rows.Close()

	if !rows.Next() {
		return domain.Source{}, domain.ErrNotFound
	}

	var (
		kind    domain.SourceType
		address string
	)

	if scanErr := rows.Scan(&kind, &address); scanErr != nil {
		return domain.Source{}, errors.Wrap(scanErr, "can't scan source row")
	}

	if rows.Err() != nil {
		return domain.Source{}, rows.Err()
	}

	return domain.Source{ID: id, Kind: kind, Address: address}, nil
}

func (repo sourceRepository) GetAll(ctx context.Context) ([]domain.Source, error) {
	const query = `
		SELECT id, kind, address
		FROM source 
		ORDER BY id ASC
	;`

	rows, rowsErr := repo.conn.Query(ctx, query)
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't exec query")
	}
	defer rows.Close()

	var (
		identifier string
		kind       domain.SourceType
		address    string
	)

	result := make([]domain.Source, 0, len(rows.RawValues()))

	for rows.Next() {
		if scanErr := rows.Scan(&identifier, &kind, &address); scanErr != nil {
			return nil, errors.Wrap(scanErr, "can't scan source row")
		}

		result = append(result, domain.Source{
			ID:      identifier,
			Kind:    kind,
			Address: address,
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}
