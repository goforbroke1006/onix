package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewMeasurementStorage creates data exchange object with db.
func NewMeasurementStorage(db *sqlx.DB) domain.MeasurementStorage {
	return &measurementRepository{db: db}
}

var _ domain.MeasurementStorage = (*measurementRepository)(nil)

type measurementRepository struct {
	db *sqlx.DB
}

func (repo measurementRepository) Store(
	ctx context.Context,
	sourceID string,
	criteriaID int64, moment time.Time, value float64,
) error {
	const query = `
		INSERT INTO measurement (source, criteria_id, moment, value, updated_at) 
		VALUES (:sourceID, :criteriaID, :moment, :value, NOW())
		ON CONFLICT (source, criteria_id, moment) 
		  DO UPDATE SET 
            value      = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
        ;
		`
	_, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"sourceID":   sourceID,
		"criteriaID": criteriaID,
		"moment":     moment.UTC().Format(time.RFC3339),
		"value":      value,
	})
	if err != nil {
		return errors.Wrap(err, "can't exec query")
	}

	return nil
}

func (repo measurementRepository) StoreBatch(
	ctx context.Context,
	sourceID string,
	criteriaID int64,
	measurements []domain.MeasurementRow,
) error {
	valuesStrs := make([]string, 0, len(measurements))
	for _, m := range measurements {
		valuesStrs = append(
			valuesStrs,
			fmt.Sprintf("('%s', %d, '%s', %f, NOW())",
				sourceID, criteriaID, m.Moment.UTC().Format(time.RFC3339), m.Value),
		)
	}

	query := fmt.Sprintf(
		`
		INSERT INTO measurement (source, criteria_id, moment, value, updated_at) 
		VALUES %s
		ON CONFLICT (source_id, criteria_id, moment) 
		  DO UPDATE SET 
			value      = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
		;	
		`,
		strings.Join(valuesStrs, ", "),
	)
	_, err := repo.db.ExecContext(ctx, query)

	return errors.Wrap(err, "can't exec query")
}
