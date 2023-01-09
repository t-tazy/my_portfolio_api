package service

import (
	"context"
	"fmt"

	"github.com/t-tazy/my_portfolio_api/auth"
	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/store"
)

type ListExercise struct {
	DB   store.Queryer
	Repo ExerciseLister
}

// contextからuser_idを取得
// 取得したユーザーの投稿を取得する
func (l *ListExercise) ListExercises(ctx context.Context) (entity.Exercises, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	exercises, err := l.Repo.ListExercises(ctx, l.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return exercises, nil
}
