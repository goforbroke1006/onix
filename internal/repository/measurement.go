package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewMeasurementRepository creates data exchange object with db.
func NewMeasurementRepository(conn *pgxpool.Pool) *measurementRepository { // nolint:revive,golint
	return &measurementRepository{
		conn: conn,
	}
}

var _ domain.MeasurementRepository = &measurementRepository{} // nolint:exhaustivestruct

type measurementRepository struct {
	conn *pgxpool.Pool
}

func (repo measurementRepository) Store(sourceID, criteriaID int64, moment time.Time, value float64) error {
	query := fmt.Sprintf(
		`
		INSERT INTO measurement (source_id, criteria_id, moment, value, updated_at) 
		VALUES (%d, %d, '%s', %f, NOW())
		ON CONFLICT (source_id, criteria_id, moment) 
		  DO UPDATE SET 
            value      = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
        ;
		`,
		sourceID, criteriaID, moment.UTC().Format(time.RFC3339), value,
	)
	_, err := repo.conn.Exec(context.TODO(), query)

	return errors.Wrap(err, "can't exec query")
}

func (repo measurementRepository) StoreBatch(sourceID, criteriaID int64, measurements []domain.MeasurementRow) error {
	valuesStrs := make([]string, 0, len(measurements))

	for _, m := range measurements {
		valuesStrs = append(
			valuesStrs,
			fmt.Sprintf("(%d, %d, '%s', %f, NOW())", sourceID, criteriaID, m.Moment.UTC().Format(time.RFC3339), m.Value),
		)
	}

	query := fmt.Sprintf(
		`
		INSERT INTO measurement (source_id, criteria_id, moment, value, updated_at) 
		VALUES %s
		ON CONFLICT (source_id, criteria_id, moment) 
		  DO UPDATE SET 
			value      = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
		;	
		`,
		strings.Join(valuesStrs, ", "),
	)
	_, err := repo.conn.Exec(context.TODO(), query)

	return errors.Wrap(err, "can't exec query")
}

func (repo measurementRepository) GetBy(
	sourceID, criteriaID int64, from, till time.Time,
) ([]domain.MeasurementRow, error) {
	query := fmt.Sprintf(`
		SELECT moment, value
		FROM measurement
		WHERE source_id = %d 
			AND criteria_id = %d 
			AND moment BETWEEN '%s' AND '%s'
		ORDER BY moment ASC
		`,
		sourceID, criteriaID, from.UTC().Format(time.RFC3339), till.UTC().Format(time.RFC3339))

	rows, err := repo.conn.Query(context.Background(), query)
	if err != nil {
		return nil, errors.Wrap(err, "can't get measurement from db")
	}
	defer rows.Close()

	result := make([]domain.MeasurementRow, 0, len(rows.RawValues()))

	var (
		moment time.Time
		value  float64
	)

	for rows.Next() {
		if err := rows.Scan(&moment, &value); err != nil {
			return nil, errors.Wrap(err, "can't scan measurement row")
		}

		result = append(result, domain.MeasurementRow{
			Moment: moment,
			Value:  value,
		})
	}

	return result, nil
}

func (repo measurementRepository) Count(sourceID, criteriaID int64, from, till time.Time) (int64, error) {
	query := fmt.Sprintf(`
	SELECT COUNT(id) 
	FROM measurement
	WHERE source_id = %d 
		AND criteria_id = %d 
		AND moment BETWEEN '%s' AND '%s'
	`, sourceID, criteriaID, from.Format(time.RFC3339), till.Format(time.RFC3339))

	rows, err := repo.conn.Query(context.Background(), query)
	if err != nil {
		return 0, errors.Wrap(err, "can't get measurements count from db")
	}
	defer rows.Close()

	if rows.Next() {
		var identifier int64
		if err := rows.Scan(&identifier); err != nil {
			return 0, errors.Wrap(err, "can't scan measurement row")
		}

		return identifier, nil
	}

	return 0, domain.ErrNotFound
}

func (repo measurementRepository) GetForPoints(
	sourceID, criteriaID int64,
	points []time.Time,
) ([]domain.MeasurementRow, error) {
	pointsInClause := make([]string, 0, len(points))

	for _, point := range points {
		timeStr := point.UTC().Format(time.RFC3339)
		pointsInClause = append(pointsInClause, fmt.Sprintf(`'%s'`, timeStr))
	}

	pointsInClauseStr := strings.Join(pointsInClause, ",")

	query := fmt.Sprintf(`
		SELECT moment, value
		FROM measurement
		WHERE source_id = %d 
			AND criteria_id = %d 
			AND moment IN (%s)
		ORDER BY moment ASC
		`,
		sourceID, criteriaID, pointsInClauseStr)

	rows, err := repo.conn.Query(context.Background(), query)
	if err != nil {
		return nil, errors.Wrap(err, "can't get measurement from db")
	}
	defer rows.Close()

	result := make([]domain.MeasurementRow, 0, len(rows.RawValues()))

	var (
		moment time.Time
		value  float64
	)

	for rows.Next() {
		if err := rows.Scan(&moment, &value); err != nil {
			return nil, errors.Wrap(err, "can't scan measurement row")
		}

		result = append(result, domain.MeasurementRow{
			Moment: moment,
			Value:  value,
		})
	}

	return result, nil
}
