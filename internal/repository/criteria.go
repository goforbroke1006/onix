package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/goforbroke1006/onix/domain"
)

// NewCriteriaRepository creates data exchange object with db
func NewCriteriaRepository(conn *pgxpool.Pool) *criteriaRepository { // nolint:revive,golint
	return &criteriaRepository{
		conn: conn,
	}
}

var _ domain.CriteriaRepository = &criteriaRepository{}

type criteriaRepository struct {
	conn *pgxpool.Pool
}

func (repo criteriaRepository) Create(
	serviceName, title string,
	selector string,
	expectedDir domain.DynamicDirType,
	interval domain.GroupingIntervalType,
) (int64, error) {
	query := fmt.Sprintf(
		`
		INSERT INTO criteria (service, title, selector, expected_dir, grouping_interval) 
		VALUES ('%s', '%s', '%s', '%s', '%s') 
		RETURNING id;`,
		serviceName, title, selector, expectedDir, interval,
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

func (repo criteriaRepository) GetAll(serviceName string) ([]domain.Criteria, error) {
	query := fmt.Sprintf(
		`
		SELECT 
			id, 
			title,
			selector, 
			expected_dir, 
			grouping_interval
		FROM criteria 
		WHERE service = '%s'
		ORDER BY id ASC
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
		title       string
		selector    string
		expectedDir domain.DynamicDirType
		interval    string
	)

	result := make([]domain.Criteria, 0, len(rows.RawValues()))

	for rows.Next() {
		if err := rows.Scan(&id, &title, &selector, &expectedDir, &interval); err != nil {
			return nil, err
		}

		duration, _ := time.ParseDuration(interval)

		result = append(result, domain.Criteria{
			ID:               id,
			Service:          serviceName,
			Title:            title,
			Selector:         selector,
			ExpectedDir:      expectedDir,
			GroupingInterval: domain.GroupingIntervalType(duration),
		})
	}

	return result, nil
}
