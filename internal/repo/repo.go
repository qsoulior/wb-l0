package repo

import "github.com/qsoulior/wb-l0/internal/entity"

type Repo interface {
	Get() ([]entity.Order, error)
	GetByID(orderID string) (*entity.Order, error)
	Create(order entity.Order) (*entity.Order, error)
}
