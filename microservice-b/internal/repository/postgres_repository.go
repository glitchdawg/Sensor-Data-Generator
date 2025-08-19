package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	
	_ "github.com/lib/pq"
	"github.com/glitchdawg/synthetic_sensors/microservice-b/internal/domain"
	sharedDomain "github.com/glitchdawg/synthetic_sensors/shared/domain"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) domain.SensorReadingRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, reading *sharedDomain.SensorReading) error {
	query := `INSERT INTO sensor_readings (id1, id2, sensor_type, value, ts) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, reading.ID1, reading.ID2, reading.SensorType, reading.Value, reading.Timestamp)
	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*sharedDomain.SensorReading, error) {
	reading := &sharedDomain.SensorReading{}
	query := `SELECT id, id1, id2, sensor_type, value, ts FROM sensor_readings WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reading.ID, &reading.ID1, &reading.ID2, &reading.SensorType, &reading.Value, &reading.Timestamp,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return reading, err
}

func (r *postgresRepository) GetByFilter(ctx context.Context, filter *sharedDomain.SensorReadingFilter) (*sharedDomain.PaginatedResponse, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	if filter.ID1 != nil {
		conditions = append(conditions, fmt.Sprintf("id1 = $%d", argCount))
		args = append(args, *filter.ID1)
		argCount++
	}
	if filter.ID2 != nil {
		conditions = append(conditions, fmt.Sprintf("id2 = $%d", argCount))
		args = append(args, *filter.ID2)
		argCount++
	}
	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("ts >= $%d", argCount))
		args = append(args, *filter.From)
		argCount++
	}
	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("ts <= $%d", argCount))
		args = append(args, *filter.To)
		argCount++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total items
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM sensor_readings %s", whereClause)
	var totalItems int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	totalPages := int(totalItems) / filter.PageSize
	if int(totalItems)%filter.PageSize > 0 {
		totalPages++
	}

	// Fetch paginated data
	query := fmt.Sprintf("SELECT id, id1, id2, sensor_type, value, ts FROM sensor_readings %s ORDER BY ts DESC LIMIT $%d OFFSET $%d", whereClause, argCount, argCount+1)
	args = append(args, filter.PageSize, offset)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []sharedDomain.SensorReading
	for rows.Next() {
		var reading sharedDomain.SensorReading
		err := rows.Scan(&reading.ID, &reading.ID1, &reading.ID2, &reading.SensorType, &reading.Value, &reading.Timestamp)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}

	return &sharedDomain.PaginatedResponse{
		Data:       readings,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r *postgresRepository) Update(ctx context.Context, id int, reading *sharedDomain.SensorReading) error {
	query := `UPDATE sensor_readings SET id1 = $1, id2 = $2, sensor_type = $3, value = $4, ts = $5 WHERE id = $6`
	result, err := r.db.ExecContext(ctx, query, reading.ID1, reading.ID2, reading.SensorType, reading.Value, reading.Timestamp, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, filter *sharedDomain.SensorReadingFilter) (int64, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	if filter.ID1 != nil {
		conditions = append(conditions, fmt.Sprintf("id1 = $%d", argCount))
		args = append(args, *filter.ID1)
		argCount++
	}
	if filter.ID2 != nil {
		conditions = append(conditions, fmt.Sprintf("id2 = $%d", argCount))
		args = append(args, *filter.ID2)
		argCount++
	}
	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("ts >= $%d", argCount))
		args = append(args, *filter.From)
		argCount++
	}
	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("ts <= $%d", argCount))
		args = append(args, *filter.To)
		argCount++
	}

	if len(conditions) == 0 {
		return 0, fmt.Errorf("at least one filter condition is required for deletion")
	}

	query := fmt.Sprintf("DELETE FROM sensor_readings WHERE %s", strings.Join(conditions, " AND "))
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}