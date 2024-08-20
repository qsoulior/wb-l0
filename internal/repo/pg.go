package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/pkg/postgres"
)

type repoPG struct {
	*postgres.Postgres
}

func NewRepoPG(pg *postgres.Postgres) Repo {
	return &repoPG{pg}
}

func (r *repoPG) Get(ctx context.Context) ([]entity.Order, error) {
	const query = "SELECT * FROM order"
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Order])
}

func (r *repoPG) GetByID(ctx context.Context, orderID string) (*entity.Order, error) {
	const query = "SELECT * FROM order WHERE id = $1"
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Order])
}

func (r *repoPG) Create(ctx context.Context, order entity.Order) (*entity.Order, error) {
	const query = "INSERT INTO order (id, data) VALUES ($1, $2) RETURNING *"
	rows, err := r.Pool.Query(ctx, query, order.OrderUID, order)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Order])
}

func (r *repoPG) CreateMany(ctx context.Context, orders []entity.Order) error {
	n := len(orders)
	if n == 0 {
		return nil
	}

	_, err := r.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"order"},
		[]string{"id", "data"},
		pgx.CopyFromSlice(n, func(i int) ([]any, error) {
			return []any{orders[i].OrderUID, orders[i]}, nil
		}))

	if err != nil {
		return err
	}

	return nil
}
