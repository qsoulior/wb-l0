package repo

import (
	"context"
	"time"

	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/pkg/cache"
)

type Cache struct {
	*cache.Cache[entity.Order]
}

func NewCache(ctx context.Context) Repo {
	cache := cache.New[entity.Order](ctx, 0, 10*time.Minute)
	return &Cache{cache}
}

func (r *Cache) Get(ctx context.Context) ([]entity.Order, error) {
	return r.Cache.Values(), nil
}

func (r *Cache) GetByID(ctx context.Context, orderID string) (*entity.Order, error) {
	order, ok := r.Cache.Get(orderID)
	if !ok {
		return nil, ErrNoRows
	}
	return &order, nil
}

func (r *Cache) Create(ctx context.Context, order entity.Order) (*entity.Order, error) {
	r.Cache.Set(order.OrderUID, order, 0)
	return &order, nil
}

func (r *Cache) CreateMany(ctx context.Context, orders []entity.Order) error {
	for _, order := range orders {
		r.Cache.Set(order.OrderUID, order, 0)
	}
	return nil
}
