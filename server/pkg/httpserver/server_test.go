package httpserver

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	want := "host:port"
	s := New(nil, "host", "port")

	got := s.server.Addr
	if got != want {
		t.Errorf("got = %v, want %v", got, want)
	}

	got1 := s.server.Handler
	if got1 != nil {
		t.Errorf("got1 = %v, want1 %v", got1, nil)
	}
}

func TestServer_Start(t *testing.T) {
	s := New(nil, "host", "port")
	want, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(want)
	s.Stop(want)
	got := s.server.BaseContext(nil)
	if want != got {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestServer_Err(t *testing.T) {
	s := New(nil, "host", "port")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	s.Stop(ctx)
	got := <-s.Err()
	if got == nil {
		t.Errorf("got = %v, want = <error>", got)
	}
}

func TestServer_Stop(t *testing.T) {
	s := New(nil, "host", "port")
	ctx := context.Background()
	s.Start(ctx)
	s.Stop(ctx)
	got := <-s.Err()
	if got == nil {
		t.Errorf("got = %v, want = <error>", got)
	}
}
