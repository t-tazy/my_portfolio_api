package service

import (
	"context"
	"fmt"

	"github.com/t-tazy/my_portfolio_api/store"
)

type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

// 引数で渡されたユーザー名のユーザーを取得し、パスワード認証を行う
// 成功すればアクセストークンを返す
func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	u, err := l.Repo.GetUser(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}
	// パスワード認証
	if err := u.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}
	jwt, err := l.TokenGenerator.GenerateToken(ctx, u)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}
	return string(jwt), nil
}
