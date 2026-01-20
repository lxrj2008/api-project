package service

import (
	"context"
	"database/sql"
	"errors"

	"liangxiong/demo/auth"
	"liangxiong/demo/dto"
	"liangxiong/demo/repository"
	"liangxiong/demo/utils"
)

// AuthService handles authentication flows.
type AuthService struct {
	repo repository.UserRepository
	jwt  *auth.JWTManager
}

// NewAuthService constructs the service.
func NewAuthService(repo repository.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{repo: repo, jwt: jwt}
}

// Login authenticates username/password.
func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByUsername(ctx, nil, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.Clone(utils.ErrUnauthorized, map[string]string{"username": "not found"}, err)
		}
		return nil, err
	}

	if !auth.VerifyPassword(user.PasswordHash, password) {
		return nil, utils.Clone(utils.ErrUnauthorized, map[string]string{"username": "invalid credentials"}, nil)
	}

	token, expiresAt, err := s.jwt.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{AccessToken: token, ExpiresAt: expiresAt}, nil
}
