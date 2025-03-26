package server

import (
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	router := http.NewServeMux()
	server := NewServer(router)

	assertEqual(t, "server address", server.Addr, ":8080")
	assertEqual(t, "read timeout", server.ReadTimeout, 10*time.Second)
	assertEqual(t, "read header timeout", server.ReadHeaderTimeout, 2*time.Second)
	assertEqual(t, "write timeout", server.WriteTimeout, 10*time.Second)
	assertEqual(t, "idle timeout", server.IdleTimeout, 60*time.Second)
}

func assertEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("expected %s %v, but got %v", name, want, got)
	}
}
