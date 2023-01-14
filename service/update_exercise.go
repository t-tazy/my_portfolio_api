package service

import (
	"context"
	"fmt"

	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/store"
)

type UpdateExercise struct {
	DB   store.Execer
	Repo ExerciseUpdater
}

func (u *UpdateExercise) UpdateExercise(ctx context.Context, id entity.ExerciseID, title, description string) error {
	e := &entity.Exercise{
		ID:          id,
		Title:       title,
		Description: description,
	}
	affected, err := u.Repo.UpdateExercise(ctx, u.DB, e)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	if *affected != 1 {
		return fmt.Errorf("failed to update: affected rows: %d", *affected)
	}
	return nil
}
