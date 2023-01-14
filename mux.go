package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/t-tazy/my_portfolio_api/auth"
	"github.com/t-tazy/my_portfolio_api/clock"
	"github.com/t-tazy/my_portfolio_api/config"
	"github.com/t-tazy/my_portfolio_api/handler"
	"github.com/t-tazy/my_portfolio_api/service"
	"github.com/t-tazy/my_portfolio_api/store"
)

// 引数としてDBへの接続情報を受け取る
// マルチプレクサを返す(ルーティング情報を持つ)
// db.Close()を実行する無名関数を返す
func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	// HTTPサーバーが稼働中か確認するための/healthエンドポイントを宣言
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析のエラーを回避するため明示的に戻り値を捨てる
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	clocker := clock.RealClocker{}
	r := store.Repository{Clocker: clocker}
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	l := &handler.Login{
		Service:   &service.Login{DB: db, Repo: &r, TokenGenerator: jwter},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)

	ae := &handler.AddExercise{
		Service:   &service.AddExercise{DB: db, Repo: &r},
		Validator: v,
	}
	le := &handler.ListExercise{
		Service: &service.ListExercise{DB: db, Repo: &r},
	}
	ge := &handler.GetExercise{
		Service: &service.GetExercise{DB: db, Repo: &r},
	}
	mux.Route("/exercises", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", ae.ServeHTTP)
		r.Get("/", le.ServeHTTP)
		r.Route("/{exerciseID}", func(r chi.Router) {
			r.Get("/", ge.ServeHTTP)
		})
	})

	// adminロールのユーザーのみがアクセス可能なエンドポイント
	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message": "admin only"}`))
		})
	})

	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	return mux, cleanup, nil
}
