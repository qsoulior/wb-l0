package service

import "github.com/qsoulior/wb-l0/internal/entity"

type Service interface {
	Get(id string) (*entity.Order, error)
	Create(order entity.Order) (*entity.Order, error)
}
