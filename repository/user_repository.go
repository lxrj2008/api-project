package repository

import (
	"context"
	"database/sql"
	"errors"

	"liangxiong/demo/model/entity"
	"liangxiong/demo/utils"
)

// sqlExecutor wraps *sql.DB or *sql.Tx to unify method calls.
type sqlExecutor struct {
	inner interface{}
}

func (s sqlExecutor) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	switch db := s.inner.(type) {
	case *sql.DB:
		return db.ExecContext(ctx, query, args...)
	case *sql.Tx:
		return db.ExecContext(ctx, query, args...)
	default:
		return nil, errors.New("unsupported executor")
	}
}

func (s sqlExecutor) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	switch db := s.inner.(type) {
	case *sql.DB:
		return db.QueryContext(ctx, query, args...)
	case *sql.Tx:
		return db.QueryContext(ctx, query, args...)
	default:
		return nil, errors.New("unsupported executor")
	}
}

func (s sqlExecutor) queryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	switch db := s.inner.(type) {
	case *sql.DB:
		return db.QueryRowContext(ctx, query, args...)
	case *sql.Tx:
		return db.QueryRowContext(ctx, query, args...)
	default:
		return nil
	}
}

// UserRepository exposes user persistence operations.
type UserRepository interface {
	GetByID(ctx context.Context, exec *sql.Tx, id string) (*entity.User, error)
	GetByUsername(ctx context.Context, exec *sql.Tx, username string) (*entity.User, error)
	List(ctx context.Context, page, size int) ([]entity.User, int64, error)
	Create(ctx context.Context, exec *sql.Tx, user *entity.User) error
	Update(ctx context.Context, exec *sql.Tx, user *entity.User) error
	Delete(ctx context.Context, exec *sql.Tx, id string) error
}

// SQLUserRepository is the concrete repository backed by SQL Server.
type SQLUserRepository struct {
	db *sql.DB
}

// NewUserRepository instantiates a repository.
func NewUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) withExecutor(exec *sql.Tx) sqlExecutor {
	if exec != nil {
		return sqlExecutor{inner: exec}
	}
	return sqlExecutor{inner: r.db}
}

// GetByID fetches a user by identifier.
func (r *SQLUserRepository) GetByID(ctx context.Context, exec *sql.Tx, id string) (*entity.User, error) {
	row := r.withExecutor(exec).queryRowContext(ctx, `SELECT id, username, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users WHERE id = @p1`, id)
	return scanUser(row)
}

// GetByUsername fetches a user by username.
func (r *SQLUserRepository) GetByUsername(ctx context.Context, exec *sql.Tx, username string) (*entity.User, error) {
	row := r.withExecutor(exec).queryRowContext(ctx, `SELECT id, username, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users WHERE username = @p1`, username)
	return scanUser(row)
}

// List returns paginated users plus total count.
func (r *SQLUserRepository) List(ctx context.Context, page, size int) ([]entity.User, int64, error) {
	exec := r.withExecutor(nil)
	offset := utils.Offset(page, size)
	limit := utils.NormalizeSize(size)
	rows, err := exec.queryContext(ctx, `SELECT id, username, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users ORDER BY created_at DESC OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY`, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int64
	row := exec.queryRowContext(ctx, `SELECT COUNT(1) FROM users`)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Create inserts a new user.
func (r *SQLUserRepository) Create(ctx context.Context, exec *sql.Tx, user *entity.User) error {
	_, err := r.withExecutor(exec).execContext(ctx, `INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, created_at, updated_at) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)`,
		user.ID, user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

// Update modifies an existing user.
func (r *SQLUserRepository) Update(ctx context.Context, exec *sql.Tx, user *entity.User) error {
	_, err := r.withExecutor(exec).execContext(ctx, `UPDATE users SET email = @p1, first_name = @p2, last_name = @p3, role = @p4, updated_at = @p5 WHERE id = @p6`,
		user.Email, user.FirstName, user.LastName, user.Role, user.UpdatedAt, user.ID)
	return err
}

// Delete removes a user row.
func (r *SQLUserRepository) Delete(ctx context.Context, exec *sql.Tx, id string) error {
	_, err := r.withExecutor(exec).execContext(ctx, `DELETE FROM users WHERE id = @p1`, id)
	return err
}

func scanUser(row *sql.Row) (*entity.User, error) {
	if row == nil {
		return nil, sql.ErrNoRows
	}
	var u entity.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
