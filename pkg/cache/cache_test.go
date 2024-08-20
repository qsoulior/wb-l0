package cache

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func emptyCache() *Cache[any] { return New[any](context.TODO(), 5*time.Minute, 10*time.Minute) }

func simpleCache() *Cache[any] {
	cache := emptyCache()
	cache.Set("key", "value", -1)
	return cache
}

func TestItem_Expired(t *testing.T) {
	tests := []struct {
		name string
		i    Item[any]
		want bool
	}{
		{"PositiveExpiredAt", Item[any]{ExpiredAt: 1}, true},
		{"ZeroExpiredAt", Item[any]{ExpiredAt: 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Expired(); got != tt.want {
				t.Errorf("Item.Expired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		expiration time.Duration
		interval   time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Cache[any]
	}{
		{"PositiveInterval", args{0, 10 * time.Minute}, &Cache[any]{items: make(map[string]Item[any]), expiration: 0, interval: 10 * time.Minute}},
		{"ZeroInterval", args{0, 0}, &Cache[any]{items: make(map[string]Item[any])}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New[any](context.TODO(), tt.args.expiration, tt.args.interval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {
	type args struct {
		key        string
		value      any
		expiration time.Duration
	}
	tests := []struct {
		name string
		c    *Cache[any]
		args args
	}{
		{"PositiveExpiration", emptyCache(), args{"key", "value", 7 * time.Minute}},
		{"ZeroExpiration", emptyCache(), args{"key", "value", 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Set(tt.args.key, tt.args.value, tt.args.expiration)
		})
	}
}

func TestCache_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *Cache[any]
		args  args
		want  any
		want1 bool
	}{
		{"NotEmptyCache", simpleCache(), args{"key"}, "value", true},
		{"EmptyCache", emptyCache(), args{"key"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCache_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Cache[any]
		args args
	}{
		{"NotEmptyCache", simpleCache(), args{"key"}},
		{"EmptyCache", emptyCache(), args{"key"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Delete(tt.args.key)
		})
	}
}

func TestCache_DeleteExpired(t *testing.T) {
	c := emptyCache()
	c.items["key"] = Item[any]{"value", -1}

	tests := []struct {
		name string
		c    *Cache[any]
	}{
		{"NotExpiredCache", simpleCache()},
		{"EmptyCache", emptyCache()},
		{"ExpiredCache", c},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.DeleteExpired()
		})
	}
}

func TestCache_gc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := simpleCache()
	c.gc(ctx)
}
