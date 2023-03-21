package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewReleaseRepository creates data exchange object with db.
func NewReleaseRepository(db *sqlx.DB) domain.ReleaseRepository {
	return &releaseRepository{db: db}
}

var _ domain.ReleaseRepository = (*releaseRepository)(nil)

type releaseRepository struct {
	db *sqlx.DB
}

func (repo releaseRepository) GetAll(serviceName string) ([]domain.Release, error) {
	const query = `
		SELECT tag, start_at 
		FROM release 
		WHERE service = :serviceName
		ORDER BY start_at;
	`

	rows, rowsErr := repo.db.NamedQueryContext(context.TODO(), query, map[string]interface{}{
		"serviceName": serviceName,
	})
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't extract releases from db")
	}
	defer func() { _ = rows.Close() }()

	var (
		releaseName string
		startAt     time.Time
	)

	var result []domain.Release

	for rows.Next() {
		if scanErr := rows.Scan(&releaseName, &startAt); scanErr != nil {
			return nil, errors.Wrap(scanErr, "can't scan release row")
		}

		result = append(result, domain.Release{
			Service: serviceName,
			Tag:     releaseName,
			StartAt: startAt,
		})
	}

	return result, nil
}

func (repo releaseRepository) Store(serviceName string, tagName string, startAt time.Time) error {
	const query = `
		INSERT INTO release (service, tag, start_at) 
		VALUES (:serviceName, :tagName, :startAt)
		ON CONFLICT (service, tag) DO UPDATE SET start_at = EXCLUDED.start_at;
	`

	_, execErr := repo.db.NamedExecContext(context.TODO(), query, map[string]interface{}{
		"serviceName": serviceName,
		"tagName":     tagName,
		"startAt":     startAt.Format(time.RFC3339),
	})

	return errors.Wrap(execErr, "can't exec query")
}

func (repo releaseRepository) GetByName(serviceName, tagName string) (*domain.Release, error) {
	const query = `
		SELECT start_at 
		FROM release 
		WHERE service = :serviceName AND tag = :tagName
		LIMIT 1;`

	rows, rowsErr := repo.db.NamedQueryContext(context.TODO(), query, map[string]interface{}{
		"serviceName": serviceName,
		"tagName":     tagName,
	})
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't get release by name from db")
	}
	defer func() { _ = rows.Close() }()

	var startAt time.Time

	if rows.Next() {
		if err := rows.Scan(&startAt); err != nil {
			return nil, errors.Wrap(err, "can't scan release row")
		}

		release := domain.Release{
			Service: serviceName,
			Tag:     tagName,
			StartAt: startAt,
		}

		return &release, nil
	}

	return nil, domain.ErrNotFound
}

func (repo releaseRepository) GetNextAfter(serviceName, tagName string) (*domain.Release, error) {
	const query = `
		SELECT tag, start_at 
		FROM release 
		WHERE service = :serviceName
		  AND start_at > (SELECT start_at FROM release WHERE service = :serviceName AND tag = :tagName)
		ORDER BY start_at
		LIMIT 1;`

	rows, rowsErr := repo.db.NamedQueryContext(context.TODO(), query, map[string]interface{}{
		"serviceName": serviceName,
		"tagName":     tagName,
	})
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't get next release from db")
	}
	defer func() { _ = rows.Close() }()

	var startAt time.Time

	if rows.Next() {
		if scanErr := rows.Scan(&tagName, &startAt); scanErr != nil {
			return nil, errors.Wrap(scanErr, "can't scan release row")
		}

		release := domain.Release{
			Service: serviceName,
			Tag:     tagName,
			StartAt: startAt,
		}

		return &release, nil
	}

	return nil, domain.ErrNotFound
}

func (repo releaseRepository) GetLast(serviceName string) (*domain.Release, error) {
	const query = `
		SELECT tag, start_at 
		FROM release 
		WHERE service = :serviceName
		ORDER BY start_at DESC
		LIMIT 1;
	`

	rows, rowsErr := repo.db.NamedQueryContext(context.TODO(), query, map[string]interface{}{
		"serviceName": serviceName,
	})
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't get least release from db")
	}
	defer func() { _ = rows.Close() }()

	var (
		tagName string
		startAt time.Time
	)

	if rows.Next() {
		if scanErr := rows.Scan(&tagName, &startAt); scanErr != nil {
			return nil, errors.Wrap(scanErr, "can't scan release row")
		}

		release := domain.Release{
			Service: serviceName,
			Tag:     tagName,
			StartAt: startAt,
		}

		return &release, nil
	}

	return nil, domain.ErrNotFound
}

func (repo releaseRepository) GetNLasts(serviceName string, count uint) ([]domain.Release, error) {
	const query = `
		SELECT tag, start_at 
		FROM release 
		WHERE service = :serviceName
		ORDER BY start_at DESC
		LIMIT :count;
	`

	rows, rowsErr := repo.db.NamedQueryContext(context.TODO(), query, map[string]interface{}{
		"serviceName": serviceName,
		"count":       count,
	})
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't get N last releases from db")
	}
	defer func() { _ = rows.Close() }()

	var (
		tagName string
		startAt time.Time
	)

	var result []domain.Release

	for rows.Next() {
		if scanErr := rows.Scan(&tagName, &startAt); scanErr != nil {
			return nil, errors.Wrap(scanErr, "can't scan release row")
		}

		result = append(result, domain.Release{
			Service: serviceName,
			Tag:     tagName,
			StartAt: startAt,
		})
	}

	return result, nil
}

func (repo releaseRepository) GetReleases(serviceName string, from, till time.Time) ([]domain.Release, error) {
	const query = `
		SELECT tag, start_at 
		FROM release 
		WHERE service = :serviceName
		  AND start_at BETWEEN :from AND :till
		ORDER BY start_at;
	`

	rows, rowsErr := repo.db.NamedQueryContext(context.TODO(), query, map[string]interface{}{
		"serviceName": serviceName,
		"from":        from.Format(time.RFC3339),
		"till":        till.Format(time.RFC3339),
	})
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't exec query")
	}
	defer func() { _ = rows.Close() }()

	var (
		tagName string
		startAt time.Time
	)

	var result []domain.Release

	for rows.Next() {
		if scanErr := rows.Scan(&tagName, &startAt); scanErr != nil {
			return nil, errors.Wrap(scanErr, "can't scan release row")
		}

		result = append(result, domain.Release{
			Service: serviceName,
			Tag:     tagName,
			StartAt: startAt,
		})
	}

	return result, nil
}
