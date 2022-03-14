package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewServiceRepository creates data exchange object with db.
func NewServiceRepository(conn *pgxpool.Pool) *serviceRepository { // nolint:revive,golint
	return &serviceRepository{
		conn: conn,
	}
}

var _ domain.ServiceRepository = &serviceRepository{} // nolint:exhaustivestruct

type serviceRepository struct {
	conn *pgxpool.Pool
}

func (repo serviceRepository) Store(title string) error {
	query := fmt.Sprintf(
		"INSERT INTO service (title) VALUES ('%s') ON CONFLICT DO NOTHING;",
		title,
	)
	_, err := repo.conn.Exec(context.TODO(), query)

	return errors.Wrap(err, "can't store service")
}

func (repo serviceRepository) GetAll() ([]domain.Service, error) {
	query := `SELECT title FROM service ORDER BY title ASC;`

	var (
		rows pgx.Rows
		err  error
	)

	if rows, err = repo.conn.Query(context.TODO(), query); err != nil {
		return nil, errors.Wrap(err, "can't exec query")
	}
	defer rows.Close()

	var title string

	result := make([]domain.Service, 0, len(rows.RawValues()))

	for rows.Next() {
		if err := rows.Scan(&title); err != nil {
			return nil, errors.Wrap(err, "can't scan service row")
		}

		result = append(result, domain.Service{
			Title: title,
		})
	}

	return result, nil
}
