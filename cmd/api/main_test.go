package main

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetEnvReturnsValue(t *testing.T) {
	t.Setenv("TEST_VALUE", "configured")

	if got := getEnv("TEST_VALUE", "fallback"); got != "configured" {
		t.Fatalf("expected configured value, got %q", got)
	}
}

func TestGetEnvReturnsFallback(t *testing.T) {
	if got := getEnv("MISSING_TEST_VALUE", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback value, got %q", got)
	}
}

func TestRunUsesInMemoryStoreAndDefaultPort(t *testing.T) {
	var gotAddr string
	err := run(
		func(string) string { return "" },
		func(addr string, h http.Handler) error {
			gotAddr = addr
			if h == nil {
				t.Fatal("expected handler")
			}
			return nil
		},
	)

	if err != nil {
		t.Fatalf("expected run to succeed: %v", err)
	}
	if gotAddr != ":8080" {
		t.Fatalf("expected default address :8080, got %q", gotAddr)
	}
}

func TestRunUsesConfiguredPort(t *testing.T) {
	getenv := func(key string) string {
		if key == "PORT" {
			return "9090"
		}
		return ""
	}

	var gotAddr string
	err := run(getenv, func(addr string, h http.Handler) error {
		gotAddr = addr
		return nil
	})

	if err != nil {
		t.Fatalf("expected run to succeed: %v", err)
	}
	if gotAddr != ":9090" {
		t.Fatalf("expected address :9090, got %q", gotAddr)
	}
}

func TestRunReturnsListenError(t *testing.T) {
	wantErr := errors.New("listen failed")

	err := run(
		func(string) string { return "" },
		func(string, http.Handler) error { return wantErr },
	)

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected listen error, got %v", err)
	}
}
