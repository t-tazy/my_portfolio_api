package service

import (
	"context"
	"fmt"

	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/store"
)

type GetExercise struct {
	DB   store.Queryer
	Repo ExerciseGetter
}

func (g *GetExercise) GetExercise(ctx context.Context, id entity.ExerciseID) (*entity.Exercise, error) {
	exercise, err := g.Repo.GetExercise(ctx, g.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	return exercise, nil
}
