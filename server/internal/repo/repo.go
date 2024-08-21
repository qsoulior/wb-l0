package repo

import (
	"context"

	"github.com/qsoulior/wb-l0/internal/entity"
)

type Repo interface {
	Get(ctx context.Context) ([]entity.Order, error)
	GetByID(ctx context.Context, orderID string) (*entity.Order, error)
	Create(ctx context.Context, order entity.Order) (*entity.Order, error)
	CreateMany(ctx context.Context, orders []entity.Order) error
}
