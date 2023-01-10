package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/t-tazy/my_portfolio_api/config"
)

// メインプロセスのrun関数を実行する
func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

// 環境変数から読み込んだポート番号でリッスンするHTTPサーバーを起動
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		return err
	}
	// メインプロセスであるrun関数の終了に合わせて
	// DBのコネクションを終了させる
	defer cleanup()
	s := NewServer(l, mux)
	return s.Run(ctx)
}
