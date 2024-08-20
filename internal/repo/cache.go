package repo

import (
	"context"
	"time"

	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/pkg/cache"
)

type repoCache struct {
	*cache.Cache[entity.Order]
}

func NewRepoCache(ctx context.Context) Repo {
	cache := cache.New[entity.Order](ctx, 0, 10*time.Minute)
	return &repoCache{cache}
}

func (r *repoCache) Get(ctx context.Context) ([]entity.Order, error) {
	return r.Cache.Values(), nil
}

func (r *repoCache) GetByID(ctx context.Context, orderID string) (*entity.Order, error) {
	order, ok := r.Cache.Get(orderID)
	if !ok {
		return nil, ErrNoRows
	}
	return &order, nil
}

func (r *repoCache) Create(ctx context.Context, order entity.Order) (*entity.Order, error) {
	r.Cache.Set(order.OrderUID, order, 0)
	return &order, nil
}

func (r *repoCache) CreateMany(ctx context.Context, orders []entity.Order) error {
	for _, order := range orders {
		r.Cache.Set(order.OrderUID, order, 0)
	}
	return nil
}
