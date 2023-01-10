package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/t-tazy/my_portfolio_api/config"
)

// httptestパッケージを使い、ServeHTTPに渡すモックを生成し、
// ルーティングが意図通りかテストする
func TestNewMux(t *testing.T) {
	t.Skip("リファクタリング中")

	ctx := context.Background()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut, Cleanup, err := NewMux(ctx, &config.Config{})
	t.Cleanup(Cleanup)
	sut.ServeHTTP(w, r)
	rsp := w.Result()
	t.Cleanup(func() { _ = rsp.Body.Close() })

	if rsp.StatusCode != http.StatusOK {
		t.Errorf("want status code 200, but %d", rsp.StatusCode)
	}
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
