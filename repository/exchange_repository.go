package repository

import (
	"context"
	"database/sql"

	"liangxiong/demo/model/entity"
	"liangxiong/demo/utils"
)

// ExchangeRepository exposes persistence operations for exchanges.
type ExchangeRepository interface {
	GetByCode(ctx context.Context, exec *sql.Tx, code string) (*entity.Exchange, error)
	List(ctx context.Context, page, size int) ([]entity.Exchange, int64, error)
	Create(ctx context.Context, exec *sql.Tx, exchange *entity.Exchange) error
	Update(ctx context.Context, exec *sql.Tx, exchange *entity.Exchange) error
	Delete(ctx context.Context, exec *sql.Tx, code string) error
}

// SQLExchangeRepository is the SQL Server implementation.
type SQLExchangeRepository struct {
	db *sql.DB
}

// NewExchangeRepository builds the repository.
func NewExchangeRepository(db *sql.DB) *SQLExchangeRepository {
	return &SQLExchangeRepository{db: db}
}

func (r *SQLExchangeRepository) withExecutor(exec *sql.Tx) sqlExecutor {
	if exec != nil {
		return sqlExecutor{inner: exec}
	}
	return sqlExecutor{inner: r.db}
}

// GetByCode retrieves a record by MQMExchangeCode.
func (r *SQLExchangeRepository) GetByCode(ctx context.Context, exec *sql.Tx, code string) (*entity.Exchange, error) {
	row := r.withExecutor(exec).queryRowContext(ctx, `SELECT MQMExchangeCode, ClearExchangeCode, GlobexExchangeCode, Description, SegType FROM TExchange WHERE MQMExchangeCode = @p1`, code)
	return scanExchange(row)
}

// List fetches paginated exchanges plus total count.
func (r *SQLExchangeRepository) List(ctx context.Context, page, size int) ([]entity.Exchange, int64, error) {
	exec := r.withExecutor(nil)
	offset := utils.Offset(page, size)
	limit := utils.NormalizeSize(size)
	rows, err := exec.queryContext(ctx, `SELECT MQMExchangeCode, ClearExchangeCode, GlobexExchangeCode, Description, SegType FROM TExchange ORDER BY MQMExchangeCode ASC OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY`, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var exchanges []entity.Exchange
	for rows.Next() {
		var e entity.Exchange
		if err := rows.Scan(&e.MQMExchangeCode, &e.ClearExchangeCode, &e.GlobexExchangeCode, &e.Description, &e.SegType); err != nil {
			return nil, 0, err
		}
		exchanges = append(exchanges, e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int64
	row := exec.queryRowContext(ctx, `SELECT COUNT(1) FROM TExchange`)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	return exchanges, total, nil
}

// Create inserts a new exchange row.
func (r *SQLExchangeRepository) Create(ctx context.Context, exec *sql.Tx, exchange *entity.Exchange) error {
	_, err := r.withExecutor(exec).execContext(ctx, `INSERT INTO TExchange (MQMExchangeCode, ClearExchangeCode, GlobexExchangeCode, Description, SegType) VALUES (@p1, @p2, @p3, @p4, @p5)`,
		exchange.MQMExchangeCode, exchange.ClearExchangeCode, exchange.GlobexExchangeCode, exchange.Description, exchange.SegType)
	return err
}

// Update modifies an existing exchange.
func (r *SQLExchangeRepository) Update(ctx context.Context, exec *sql.Tx, exchange *entity.Exchange) error {
	_, err := r.withExecutor(exec).execContext(ctx, `UPDATE TExchange SET ClearExchangeCode = @p1, GlobexExchangeCode = @p2, Description = @p3, SegType = @p4 WHERE MQMExchangeCode = @p5`,
		exchange.ClearExchangeCode, exchange.GlobexExchangeCode, exchange.Description, exchange.SegType, exchange.MQMExchangeCode)
	return err
}

// Delete removes a record by MQMExchangeCode.
func (r *SQLExchangeRepository) Delete(ctx context.Context, exec *sql.Tx, code string) error {
	_, err := r.withExecutor(exec).execContext(ctx, `DELETE FROM TExchange WHERE MQMExchangeCode = @p1`, code)
	return err
}

func scanExchange(row *sql.Row) (*entity.Exchange, error) {
	if row == nil {
		return nil, sql.ErrNoRows
	}
	var e entity.Exchange
	if err := row.Scan(&e.MQMExchangeCode, &e.ClearExchangeCode, &e.GlobexExchangeCode, &e.Description, &e.SegType); err != nil {
		return nil, err
	}
	return &e, nil
}
