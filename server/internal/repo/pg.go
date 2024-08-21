package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/pkg/postgres"
)

type PG struct {
	*postgres.Postgres
}

func NewPG(pg *postgres.Postgres) Repo {
	return &PG{pg}
}

type Row struct {
	ID   string
	Data entity.Order
}

func rowToStruct(r pgx.CollectableRow) (entity.Order, error) {
	row, err := pgx.RowToStructByPos[Row](r)
	if err != nil {
		return entity.Order{}, err
	}
	return row.Data, err
}

func (r *PG) Get(ctx context.Context) ([]entity.Order, error) {
	const query = `SELECT * FROM "order"`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, rowToStruct)
}

func (r *PG) GetByID(ctx context.Context, orderID string) (*entity.Order, error) {
	const query = `SELECT * FROM "order" WHERE id = $1`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	order, err := pgx.CollectExactlyOneRow(rows, rowToStruct)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, ErrNoRows
	case errors.Is(err, pgx.ErrTooManyRows):
		return nil, ErrTooManyRows
	case err != nil:
		return nil, err
	default:
		return &order, nil
	}
}

func (r *PG) Create(ctx context.Context, order entity.Order) (*entity.Order, error) {
	const query = `INSERT INTO "order" (id, data) VALUES ($1, $2) RETURNING *`
	rows, err := r.Pool.Query(ctx, query, order.OrderUID, order)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectExactlyOneRow(rows, rowToStruct)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PG) CreateMany(ctx context.Context, orders []entity.Order) error {
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
