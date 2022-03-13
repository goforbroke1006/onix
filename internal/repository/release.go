package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/goforbroke1006/onix/domain"
)

// NewReleaseRepository creates data exchange object with db
func NewReleaseRepository(conn *pgxpool.Pool) *releaseRepository { // nolint:golint
	return &releaseRepository{
		conn: conn,
	}
}

var (
	_ domain.ReleaseRepository = &releaseRepository{}
)

type releaseRepository struct {
	conn *pgxpool.Pool
}

func (repo releaseRepository) GetAll(serviceName string) ([]domain.Release, error) {
	query := fmt.Sprintf(
		`
		SELECT id, name, start_at 
		FROM release 
		WHERE service = '%s'
		ORDER BY start_at ASC
		;`,
		serviceName,
	)
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id          int64
		releaseName string
		startAt     time.Time
	)

	result := make([]domain.Release, 0, len(rows.RawValues()))
	for rows.Next() {
		if err := rows.Scan(&id, &releaseName, &startAt); err != nil {
			return nil, err
		}
		result = append(result, domain.Release{
			ID:      id,
			Service: serviceName,
			Name:    releaseName,
			StartAt: startAt,
		})
	}

	return result, nil
}

func (repo releaseRepository) Store(serviceName string, releaseName string, startAt time.Time) error {
	query := fmt.Sprintf(
		"INSERT INTO release (service, name, start_at) VALUES ('%s', '%s', '%s');",
		serviceName, releaseName, startAt.Format(time.RFC3339),
	)
	_, err := repo.conn.Exec(context.TODO(), query)
	return err
}

func (repo releaseRepository) GetByName(serviceName, releaseName string) (*domain.Release, error) {
	query := fmt.Sprintf(
		`
		SELECT id, name, start_at 
		FROM release 
		WHERE service = '%s' AND name = '%s'
		LIMIT 1
		;`,
		serviceName, releaseName,
	)
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id      int64
		startAt time.Time
	)

	if rows.Next() {
		if err := rows.Scan(&id, &releaseName, &startAt); err != nil {
			return nil, err
		}

		release := domain.Release{
			ID:      id,
			Service: serviceName,
			Name:    releaseName,
			StartAt: startAt,
		}

		return &release, nil
	}

	return nil, domain.ErrNotFound
}

func (repo releaseRepository) GetNextAfter(serviceName, releaseName string) (*domain.Release, error) {
	query := fmt.Sprintf(
		`
		SELECT id, name, start_at 
		FROM release 
		WHERE service = '%s'
		  AND start_at > (SELECT start_at FROM release WHERE service = '%s' AND name = '%s')
		ORDER BY start_at ASC
		LIMIT 1
		;`,
		serviceName,
		serviceName, releaseName,
	)
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		identifier int64
		startAt    time.Time
	)

	if rows.Next() {
		if err := rows.Scan(&identifier, &releaseName, &startAt); err != nil {
			return nil, err
		}

		release := domain.Release{
			ID:      identifier,
			Service: serviceName,
			Name:    releaseName,
			StartAt: startAt,
		}

		return &release, nil
	}

	return nil, domain.ErrNotFound
}

func (repo releaseRepository) GetLast(serviceName string) (*domain.Release, error) {
	query := fmt.Sprintf(
		`
SELECT id, name, start_at 
FROM release 
WHERE service = '%s'
ORDER BY start_at DESC
LIMIT 1
;`,
		serviceName,
	)
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id          int64
		releaseName string
		startAt     time.Time
	)

	if rows.Next() {
		if err := rows.Scan(&id, &releaseName, &startAt); err != nil {
			return nil, err
		}

		release := domain.Release{
			ID:      id,
			Service: serviceName,
			Name:    releaseName,
			StartAt: startAt,
		}

		return &release, nil
	}

	return nil, domain.ErrNotFound
}

func (repo releaseRepository) GetNLasts(serviceName string, count uint) ([]domain.Release, error) {
	query := fmt.Sprintf(
		`
		SELECT id, name, start_at 
		FROM release 
		WHERE service = '%s'
		ORDER BY start_at DESC
		LIMIT %d
		;`,
		serviceName, count,
	)
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id          int64
		releaseName string
		startAt     time.Time
	)

	result := make([]domain.Release, 0, len(rows.RawValues()))
	for rows.Next() {
		if err := rows.Scan(&id, &releaseName, &startAt); err != nil {
			return nil, err
		}
		result = append(result, domain.Release{
			ID:      id,
			Service: serviceName,
			Name:    releaseName,
			StartAt: startAt,
		})
	}

	return result, nil
}

func (repo releaseRepository) GetReleases(serviceName string, from, till time.Time) ([]domain.Release, error) {
	query := fmt.Sprintf(
		`
SELECT id, name, start_at 
FROM release 
WHERE service = '%s' 
  AND start_at BETWEEN '%s' AND '%s'
ORDER BY start_at ASC
;`,
		serviceName, from.Format(time.RFC3339), till.Format(time.RFC3339),
	)
	rows, err := repo.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id          int64
		releaseName string
		startAt     time.Time
	)

	result := make([]domain.Release, 0, len(rows.RawValues()))
	for rows.Next() {
		if err := rows.Scan(&id, &releaseName, &startAt); err != nil {
			return nil, err
		}
		result = append(result, domain.Release{
			ID:      id,
			Service: serviceName,
			Name:    releaseName,
			StartAt: startAt,
		})
	}

	return result, nil
}
