package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/internal/repo"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := repo.NewMockRepo(ctrl)
	cache := repo.NewMockRepo(ctrl)

	want := &service{db, cache}
	if got := New(db, cache); !reflect.DeepEqual(got, want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func Test_service_Init(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(db *repo.MockRepo, cache *repo.MockRepo)
		wantErr bool
	}{
		{"WithoutErr", func(db, cache *repo.MockRepo) {
			db.EXPECT().Get(gomock.Any()).Return([]entity.Order{}, nil)
			cache.EXPECT().CreateMany(gomock.Any(), gomock.Eq([]entity.Order{})).Return(nil)
		}, false},
		{"WithErr", func(db, cache *repo.MockRepo) {
			db.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf(""))
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := repo.NewMockRepo(ctrl)
			cache := repo.NewMockRepo(ctrl)
			tt.prepare(db, cache)
			s := &service{db, cache}

			if err := s.Init(nil); (err != nil) != tt.wantErr {
				t.Errorf("service.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_Get(t *testing.T) {
	type args struct {
		orderID string
	}
	tests := []struct {
		name    string
		prepare func(cache *repo.MockRepo)
		args    args
		want    *entity.Order
		wantErr bool
	}{
		{"WithoutErr", func(cache *repo.MockRepo) {
			cache.EXPECT().GetByID(gomock.Any(), gomock.Eq("1")).Return(new(entity.Order), nil)
		}, args{"1"}, new(entity.Order), false},
		{"WithErr", func(cache *repo.MockRepo) {
			cache.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf(""))
		}, args{""}, nil, true},
		{"WithErrNoRows", func(cache *repo.MockRepo) {
			cache.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, repo.ErrNoRows)
		}, args{""}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cache := repo.NewMockRepo(ctrl)
			tt.prepare(cache)
			s := &service{nil, cache}

			got, err := s.Get(nil, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Create(t *testing.T) {
	type args struct {
		order entity.Order
	}
	tests := []struct {
		name    string
		prepare func(db *repo.MockRepo, cache *repo.MockRepo)
		args    args
		want    *entity.Order
		wantErr bool
	}{
		{"WithGetByIdErrNil", func(db, cache *repo.MockRepo) {
			cache.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil)
		}, args{}, nil, true},
		{"WithGetByIdErr", func(db, cache *repo.MockRepo) {
			cache.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf(""))
		}, args{}, nil, true},
		{"WithDbCreateErr", func(db, cache *repo.MockRepo) {
			cache.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, repo.ErrNoRows)
			db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf(""))
		}, args{}, nil, true},
		{"WithCacheCreateErr", func(db, cache *repo.MockRepo) {
			cache.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, repo.ErrNoRows)
			db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil)
			cache.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf(""))
		}, args{}, nil, true},
		{"WithoutErr", func(db, cache *repo.MockRepo) {
			cache.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, repo.ErrNoRows)
			db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil)
			cache.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil)
		}, args{}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			db := repo.NewMockRepo(ctrl)
			cache := repo.NewMockRepo(ctrl)
			tt.prepare(db, cache)
			s := &service{db, cache}

			got, err := s.Create(nil, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
