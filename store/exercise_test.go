package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/t-tazy/my_portfolio_api/clock"
	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/testutil"
	"github.com/t-tazy/my_portfolio_api/testutil/fixture"
)

// 実際のRDBMSを使ってテストする
func TestRepository_ListExercises(t *testing.T) {
	ctx := context.Background()
	// entity.Exerciseを作成する他のテストケースと混ざるとfail
	// そのため、トランザクションをはることでこのテストケースの中だけのテーブル状態にする
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// このテストケースが完了したらもとに戻す
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	wantUserID, wants := prepareExercises(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListExercises(ctx, tx, wantUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(gots, wants); len(diff) != 0 {
		t.Errorf("differs: (-got +want)\n%s", diff)
	}
}

// ダミーユーザーを作成し、保存するテストヘルパー
func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()

	u := fixture.User(nil) // ダミーユーザー作成
	result, err := db.ExecContext(ctx,
		`INSERT INTO users (name, password, role, created, modified)
		VALUES (?, ?, ?, ?, ?);`,
		u.Name, u.Password, u.Role, u.Created, u.Modified,
	)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("got user_id: %v", err)
	}
	return entity.UserID(id)
}

// exercisesテーブルの状態を整えるヘルパー
// exercisesテーブルにダミーデータを保存し、期待値を返す
func prepareExercises(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Exercises) {
	t.Helper()

	userID := prepareUser(ctx, t, con)
	otherUserID := prepareUser(ctx, t, con)
	// 一度データをきれいにする
	if _, err := con.ExecContext(ctx, "DELETE FROM exercises;"); err != nil {
		t.Logf("failed to initialize exercises: %v", err)
	}

	c := clock.FixedClocker{}
	wants := entity.Exercises{
		{
			UserID:      userID,
			Title:       "exercise 1",
			Description: "want exercise 1",
			Created:     c.Now(),
			Modified:    c.Now(),
		},
		{
			UserID:      userID,
			Title:       "exercise 2",
			Description: "want exercise 2",
			Created:     c.Now(),
			Modified:    c.Now(),
		},
	}
	exercises := entity.Exercises{
		wants[0],
		{
			UserID:      otherUserID,
			Title:       "not want exercise",
			Description: "",
			Created:     c.Now(),
			Modified:    c.Now(),
		},
		wants[1],
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO exercises (user_id, title, description, created, modified)
		VALUES
			(?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?);`,
		exercises[0].UserID, exercises[0].Title, exercises[0].Description, exercises[0].Created, exercises[0].Modified,
		exercises[1].UserID, exercises[1].Title, exercises[1].Description, exercises[1].Created, exercises[1].Modified,
		exercises[2].UserID, exercises[2].Title, exercises[2].Description, exercises[2].Created, exercises[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	// 一つのINSERT文で複数のレコードを作成した場合のLastInsertIdメソッドの戻り値は
	// MySQLでは1つ目のレコードのIDになる
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	exercises[0].ID = entity.ExerciseID(id)
	exercises[1].ID = entity.ExerciseID(id + 1)
	exercises[2].ID = entity.ExerciseID(id + 2)
	return userID, wants
}

// mockを使ってテストする
func TestRepository_AddExercise(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okExercise := &entity.Exercise{
		UserID:      30,
		Title:       "ok Exericse",
		Description: "test exercise",
		Created:     c.Now(),
		Modified:    c.Now(),
	}

	// mockデータベース接続と空のmockを生成
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	// mock化
	mock.ExpectExec(
		// エスケープが必要
		`INSERT INTO exercises \(user_id, title, description, created, modified\) VALUES \(\?, \?, \?, \?, \?\)`,
	).WithArgs(okExercise.UserID, okExercise.Title, okExercise.Description, okExercise.Created, okExercise.Modified).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddExercise(ctx, xdb, okExercise); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
