package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewServiceRepository creates data exchange object with db.
func NewServiceRepository(conn *pgxpool.Pool) domain.ServiceRepository {
	return &serviceRepository{conn: conn}
}

var _ domain.ServiceRepository = (*serviceRepository)(nil)

type serviceRepository struct {
	conn *pgxpool.Pool
}

func (repo serviceRepository) Store(ctx context.Context, id string) error {
	const query = "INSERT INTO service (id) VALUES (:id) ON CONFLICT DO NOTHING;"

	_, err := repo.conn.Exec(ctx, query, id)

	return errors.Wrap(err, "can't store service")
}

func (repo serviceRepository) GetAll(ctx context.Context) ([]domain.Service, error) {
	const query = `SELECT id FROM service ORDER BY id ASC;`

	rows, rowsErr := repo.conn.Query(ctx, query)
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't exec query")
	}
	defer rows.Close()

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
