package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/example/go-api/auth"
	"github.com/example/go-api/dto"
	"github.com/example/go-api/model/entity"
	"github.com/example/go-api/repository"
	"github.com/example/go-api/utils"
)

// UserService exposes application use cases for users.
type UserService struct {
	repo repository.UserRepository
	db   *sql.DB
}

// NewUserService constructs the service.
func NewUserService(db *sql.DB, repo repository.UserRepository) *UserService {
	return &UserService{db: db, repo: repo}
}

// ListUsers returns paginated list.
func (s *UserService) ListUsers(ctx context.Context, page, size int) (*dto.UserListResponse, error) {
	users, total, err := s.repo.List(ctx, page, size)
	if err != nil {
		return nil, err
	}

	items := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		items = append(items, mapUserToDTO(&u))
	}

	return &dto.UserListResponse{Total: total, Items: items}, nil
}

// GetUser fetches a single user.
func (s *UserService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, nil, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.Clone(utils.ErrNotFound, map[string]string{"id": id}, err)
		}
		return nil, err
	}
	resp := mapUserToDTO(user)
	return &resp, nil
}

// CreateUser inserts a new user row.
func (s *UserService) CreateUser(ctx context.Context, req dto.UserCreateRequest) (*dto.UserResponse, error) {
	if _, err := s.repo.GetByUsername(ctx, nil, req.Username); err == nil {
		return nil, utils.Clone(utils.ErrBadRequest, map[string]string{"username": "already exists"}, nil)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &entity.User{
		ID:           utils.NewID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.repo.Create(ctx, tx, user); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp := mapUserToDTO(user)
	return &resp, nil
}

// UpdateUser updates selected fields.
func (s *UserService) UpdateUser(ctx context.Context, id string, req dto.UserUpdateRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, nil, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.Clone(utils.ErrNotFound, map[string]string{"id": id}, err)
		}
		return nil, err
	}

	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Role = req.Role
	user.UpdatedAt = time.Now().UTC()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.repo.Update(ctx, tx, user); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp := mapUserToDTO(user)
	return &resp, nil
}

// DeleteUser deletes by id.
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if _, err := s.repo.GetByID(ctx, nil, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.Clone(utils.ErrNotFound, map[string]string{"id": id}, err)
		}
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func mapUserToDTO(u *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
