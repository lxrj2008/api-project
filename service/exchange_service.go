package service

import (
	"context"
	"database/sql"
	"errors"

	"liangxiong/demo/dto"
	"liangxiong/demo/model/entity"
	"liangxiong/demo/repository"
	"liangxiong/demo/utils"
)

// ExchangeService orchestrates exchange workflows.
type ExchangeService struct {
	repo repository.ExchangeRepository
	db   *sql.DB
}

// NewExchangeService creates a service.
func NewExchangeService(db *sql.DB, repo repository.ExchangeRepository) *ExchangeService {
	return &ExchangeService{repo: repo, db: db}
}

// ListExchanges returns paginated exchanges.
func (s *ExchangeService) ListExchanges(ctx context.Context, page, size int) (*dto.ExchangeListResponse, error) {
	exchanges, total, err := s.repo.List(ctx, page, size)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ExchangeResponse, 0, len(exchanges))
	for _, exch := range exchanges {
		items = append(items, mapExchangeToDTO(&exch))
	}

	return &dto.ExchangeListResponse{Total: total, Items: items}, nil
}

// GetExchange fetches a single exchange.
func (s *ExchangeService) GetExchange(ctx context.Context, code string) (*dto.ExchangeResponse, error) {
	exchange, err := s.repo.GetByCode(ctx, nil, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.Clone(utils.ErrNotFound, map[string]string{"mqmExchangeCode": code}, err)
		}
		return nil, err
	}
	resp := mapExchangeToDTO(exchange)
	return &resp, nil
}

// CreateExchange inserts a new row.
func (s *ExchangeService) CreateExchange(ctx context.Context, req dto.ExchangeCreateRequest) (*dto.ExchangeResponse, error) {
	if _, err := s.repo.GetByCode(ctx, nil, req.MQMExchangeCode); err == nil {
		return nil, utils.Clone(utils.ErrBadRequest, map[string]string{"mqmExchangeCode": "already exists"}, nil)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	exchange := &entity.Exchange{
		MQMExchangeCode:    req.MQMExchangeCode,
		ClearExchangeCode:  req.ClearExchangeCode,
		GlobexExchangeCode: req.GlobexExchangeCode,
		Description:        toNullString(req.Description),
		SegType:            req.SegType,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.repo.Create(ctx, tx, exchange); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp := mapExchangeToDTO(exchange)
	return &resp, nil
}

// UpdateExchange modifies a row by code.
func (s *ExchangeService) UpdateExchange(ctx context.Context, code string, req dto.ExchangeUpdateRequest) (*dto.ExchangeResponse, error) {
	exchange, err := s.repo.GetByCode(ctx, nil, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.Clone(utils.ErrNotFound, map[string]string{"mqmExchangeCode": code}, err)
		}
		return nil, err
	}

	exchange.ClearExchangeCode = req.ClearExchangeCode
	exchange.GlobexExchangeCode = req.GlobexExchangeCode
	exchange.Description = toNullString(req.Description)
	exchange.SegType = req.SegType

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.repo.Update(ctx, tx, exchange); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp := mapExchangeToDTO(exchange)
	return &resp, nil
}

// DeleteExchange removes a row by code.
func (s *ExchangeService) DeleteExchange(ctx context.Context, code string) error {
	if _, err := s.repo.GetByCode(ctx, nil, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.Clone(utils.ErrNotFound, map[string]string{"mqmExchangeCode": code}, err)
		}
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.Delete(ctx, tx, code); err != nil {
		return err
	}

	return tx.Commit()
}

func mapExchangeToDTO(e *entity.Exchange) dto.ExchangeResponse {
	return dto.ExchangeResponse{
		MQMExchangeCode:    e.MQMExchangeCode,
		ClearExchangeCode:  e.ClearExchangeCode,
		GlobexExchangeCode: e.GlobexExchangeCode,
		Description:        nullStringToPtr(e.Description),
		SegType:            e.SegType,
	}
}

func toNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func nullStringToPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	v := value.String
	return &v
}
