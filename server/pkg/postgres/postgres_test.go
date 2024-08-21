package postgres

import (
	"context"
	"testing"
)

func TestPostgres_New(t *testing.T) {
	got, got1 := New(context.Background(), "")

	if got != nil {
		t.Errorf("got = %v, want %v", got, nil)
	}

	if got1 == nil {
		t.Errorf("got1 = %v, want1 = <error>", got)
	}
}

func TestPostgres_Close(t *testing.T) {
	p := &Postgres{}
	p.Close()
}
