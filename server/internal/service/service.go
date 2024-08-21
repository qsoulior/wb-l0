package service

import (
	"context"
	"errors"

	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/internal/repo"
)

type Service interface {
	Init(ctx context.Context) error
	Get(ctx context.Context, orderID string) (*entity.Order, error)
	Create(ctx context.Context, order entity.Order) (*entity.Order, error)
}

type service struct {
	db    repo.Repo
	cache repo.Repo
}

func New(db, cache repo.Repo) Service { return &service{db, cache} }

// Init moves orders from database to cache.
func (s *service) Init(ctx context.Context) error {
	orders, err := s.db.Get(ctx)
	if err != nil {
		return err
	}

	return s.cache.CreateMany(ctx, orders)
}

// Get returns order by its id from cache.
func (s *service) Get(ctx context.Context, orderID string) (*entity.Order, error) {
	order, err := s.cache.GetByID(ctx, orderID)
	if errors.Is(err, repo.ErrNoRows) || errors.Is(err, repo.ErrTooManyRows) {
		return nil, ErrNotExist
	} else if err != nil {
		return nil, err
	}

	return order, nil
}

// Create creates order in database and cache and returns it.
func (s *service) Create(ctx context.Context, order entity.Order) (*entity.Order, error) {
	_, err := s.cache.GetByID(ctx, order.OrderUID)
	if err == nil {
		return nil, ErrExists
	}

	if !errors.Is(err, repo.ErrNoRows) {
		return nil, err
	}

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
