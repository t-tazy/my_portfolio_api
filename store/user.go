package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/t-tazy/my_portfolio_api/entity"
)

// ユーザーを登録する
// 引数として受け取った*entity.User型のIDフィールドを更新することで
// 呼び出し元にRDBMSより発行されたIDを伝える
func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	sql := `INSERT INTO users
			(name, password, role, created, modified)
			VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(
		ctx, sql, u.Name, u.Password, u.Role, u.Created, u.Modified,
	)
	if err != nil {
		// MySQLに起因するエラー
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
		}
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = entity.UserID(id)
	return nil
}

// ユーザーを取得する
// usersテーブルのnameカラムで絞り込む。nameカラムはuniqueである。
func (r *Repository) GetUser(
	ctx context.Context, db Queryer, name string,
) (*entity.User, error) {
	u := &entity.User{}
	sql := `SELECT
			id, name, password, role, created, modified
			FROM users WHERE name = ?`
	if err := db.GetContext(ctx, u, sql, name); err != nil {
		return nil, err
	}
	return u, nil
}
