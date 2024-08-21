package service

import (
	"context"
	"errors"

	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/internal/repo"
)

type Service interface {
	Get(ctx context.Context, orderID string) (*entity.Order, error)
	Create(ctx context.Context, order entity.Order) (*entity.Order, error)
}

type service struct {
	db    repo.Repo
	cache repo.Repo
}

func New(db, cache repo.Repo) Service { return &service{db, cache} }

func (s *service) Get(ctx context.Context, orderID string) (*entity.Order, error) {
	order, err := s.db.GetByID(ctx, orderID)
	if errors.Is(err, repo.ErrNoRows) || errors.Is(err, repo.ErrTooManyRows) {
		return nil, ErrNotExist
	} else if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) Create(ctx context.Context, order entity.Order) (*entity.Order, error) {
	o, err := s.db.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	_, err = s.cache.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return o, nil
}
