package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 実DBを使ったテスト用のヘルパー
// テスト実行環境によって接続先を切り替える
func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	port := 33306
	// 環境変数CIはGithub Actions上しか定義されていないことを想定
	if _, defined := os.LookupEnv("CI"); defined {
		port = 3306
	}
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("portfolio:portfolio@tcp(127.0.0.1:%d)/portfolio?parseTime=true", port),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return sqlx.NewDb(db, "mysql")
}
