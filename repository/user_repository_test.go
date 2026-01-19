package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/example/go-api/model/entity"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	expected := entity.User{
		ID:           "user-1",
		Username:     "alice",
		Email:        "alice@example.com",
		PasswordHash: "hash",
		FirstName:    "Alice",
		LastName:     "Lee",
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "first_name", "last_name", "role", "created_at", "updated_at"}).
		AddRow(expected.ID, expected.Username, expected.Email, expected.PasswordHash, expected.FirstName, expected.LastName, expected.Role, expected.CreatedAt, expected.UpdatedAt)

	mock.ExpectQuery(`SELECT id, username, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users WHERE id = @p1`).
		WithArgs(expected.ID).
		WillReturnRows(mockRows)

	user, err := repo.GetByID(context.Background(), nil, expected.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Username != expected.Username {
		t.Fatalf("expected %s got %s", expected.Username, user.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := &entity.User{
		ID:           "user-1",
		Username:     "alice",
		Email:        "alice@example.com",
		PasswordHash: "hash",
		FirstName:    "Alice",
		LastName:     "Lee",
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, created_at, updated_at) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)`)).
		WithArgs(user.ID, user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	if err := repo.Create(context.Background(), tx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}
