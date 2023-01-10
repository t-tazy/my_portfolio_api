package service

import (
	"context"
	"fmt"

	"github.com/t-tazy/my_portfolio_api/auth"
	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/store"
)

type AddExercise struct {
	DB   store.Execer
	Repo ExerciseAdder
}

// *entity.Exercise型を引数で受け取った値で初期化し、RDBMSへ保存
// user_idはcontextから習得する
func (a *AddExercise) AddExercise(ctx context.Context, title, description string) (*entity.Exercise, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	e := &entity.Exercise{
		UserID:      id,
		Title:       title,
		Description: description,
	}
	err := a.Repo.AddExercise(ctx, a.DB, e)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return e, nil
}
