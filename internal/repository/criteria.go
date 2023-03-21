package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewCriteriaRepository creates data exchange object with db.
func NewCriteriaRepository(db *sqlx.DB) domain.CriteriaRepository {
	return &criteriaRepository{db: db}
}

var _ domain.CriteriaRepository = (*criteriaRepository)(nil)

type criteriaRepository struct {
	db *sqlx.DB
}

func (repo criteriaRepository) Create(
	ctx context.Context,
	serviceName, title string,
	selector string,
	direction domain.DirectionType,
	interval domain.GroupingIntervalType,
) (int64, error) {
	const query = `
		INSERT INTO criteria (service, title, selector, direction, "interval") 
		VALUES (:serviceName, :title, :selector, :direction, :interval) 
		RETURNING id;`

	rows, rowsErr := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"serviceName": serviceName,
		"title":       title,
		"selector":    selector,
		"direction":   direction,
		"interval":    interval,
	})
	if rowsErr != nil {
		return 0, errors.Wrap(rowsErr, "can't store criteria in db")
	}
	defer func() { _ = rows.Close() }()

	var identifier int64
	if rows.Next() {
		if scanErr := rows.Scan(&identifier); scanErr != nil {
			return 0, errors.Wrap(scanErr, "can't scan criteria row")
		}

		return identifier, nil
	}

	if rows.Err() != nil {
		return 0, rows.Err()
	}

	return 0, domain.ErrNotFound
}

func (repo criteriaRepository) GetAll(ctx context.Context, serviceName string) ([]domain.Criteria, error) {
	const query = `
		SELECT 
			id, 
			title,
			selector, 
			direction, 
			"interval"
		FROM criteria 
		WHERE service = :serviceName
		ORDER BY id;`

	rows, rowsErr := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"serviceName": serviceName,
	})
	if rowsErr != nil {
		return nil, errors.Wrap(rowsErr, "can't extract criteria from db")
	}
	defer func() { _ = rows.Close() }()

	var (
		identifier  int64
		title       string
		selector    string
		expectedDir domain.DirectionType
		interval    string
	)

	var result []domain.Criteria

	for rows.Next() {
		if scanErr := rows.Scan(&identifier, &title, &selector, &expectedDir, &interval); scanErr != nil {
			return nil, errors.Wrap(scanErr, "can't scan criteria row")
		}

		duration, _ := time.ParseDuration(interval)

		result = append(result, domain.Criteria{
			ID:        identifier,
			Service:   serviceName,
			Title:     title,
			Selector:  selector,
			Direction: expectedDir,
			Interval:  domain.GroupingIntervalType(duration),
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (repo criteriaRepository) GetByID(ctx context.Context, identifier int64) (domain.Criteria, error) {
	const query = `
		SELECT service, title, selector, direction, "interval"
		FROM criteria 
		WHERE id = :id
	`

	rows, rowsErr := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"id": identifier,
	})
	if rowsErr != nil {
		return domain.Criteria{}, errors.Wrap(rowsErr, "can't exec query")
	}
	defer func() { _ = rows.Close() }()

	var (
		service   string
		title     string
		selector  string
		direction domain.DirectionType
		interval  string
	)

	if rows.Next() {
		if scanErr := rows.Scan(&service, &title, &selector, &direction, interval); scanErr != nil {
			return domain.Criteria{}, errors.Wrap(scanErr, "can't scan criteria row")
		}

		release := domain.Criteria{
			ID:        identifier,
			Service:   service,
			Title:     title,
			Selector:  selector,
			Direction: direction,
			Interval:  domain.MustParseGroupingIntervalType(interval),
		}

		return release, nil
	}

	if rows.Err() != nil {
		return domain.Criteria{}, rows.Err()
	}

	return domain.Criteria{}, domain.ErrNotFound
}
