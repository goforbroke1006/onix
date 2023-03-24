package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewServiceStorage creates data exchange object with db.
func NewServiceStorage(db *sqlx.DB) domain.ServiceStorage {
	return &serviceRepository{db: db}
}

var _ domain.ServiceStorage = (*serviceRepository)(nil)

type serviceRepository struct {
	db *sqlx.DB
}

func (repo serviceRepository) Store(ctx context.Context, id string) error {
	const query = `
		INSERT INTO service (id) 
		VALUES (:id) 
		ON CONFLICT DO NOTHING;
	`

	_, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id": id,
	})

	return errors.Wrap(err, "can't store service")
}

func (repo serviceRepository) GetAll(ctx context.Context) ([]domain.Service, error) {
	const query = `SELECT id FROM service ORDER BY id ASC;`

	rows, rowsErr := repo.db.QueryContext(ctx, query)
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't exec query")
	}
	defer func() { _ = rows.Close() }()

	var id string

	var result []domain.Service

	for rows.Next() {
		if scanErr := rows.Scan(&id); scanErr != nil {
			return nil, errors.Wrap(scanErr, "can't scan service row")
		}

		result = append(result, domain.Service{ID: id})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}
