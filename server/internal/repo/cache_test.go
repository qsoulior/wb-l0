package repo

import (
	"context"
	"reflect"
	"testing"

	"github.com/qsoulior/wb-l0/internal/entity"
)

func TestNewCache(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if got := NewCache(ctx); got == nil {
		t.Errorf("NewCache() = %v, want <not nil>", got)
	}
}

func TestCache_Get(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(r Repo)
		want    []entity.Order
	}{
		{"EmptyCache", func(r Repo) {}, make([]entity.Order, 0)},
		{"NotEmptyCache", func(r Repo) { r.Create(nil, entity.Order{}) }, make([]entity.Order, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			r := NewCache(ctx)
			tt.prepare(r)

			got, _ := r.Get(nil)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_GetByID(t *testing.T) {
	type args struct {
		orderID string
	}
	tests := []struct {
		name    string
		prepare func(r Repo)
		args    args
		want    *entity.Order
		wantErr bool
	}{
		{"OrderExists", func(r Repo) { r.Create(nil, entity.Order{OrderUID: "1"}) }, args{"1"}, &entity.Order{OrderUID: "1"}, false},
		{"OrderDoesNotExist", func(r Repo) {}, args{"1"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			r := NewCache(ctx)
			tt.prepare(r)

			got, err := r.GetByID(nil, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cache.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := NewCache(ctx)
	want := &entity.Order{OrderUID: "1"}
	r.Create(ctx, *want)
	got, _ := r.GetByID(ctx, "1")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCache_CreateMany(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := NewCache(ctx)
	want := []entity.Order{{OrderUID: "1"}, {OrderUID: "2"}}
	r.CreateMany(ctx, want)
	got, _ := r.Get(ctx)
	if len(got) != len(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
