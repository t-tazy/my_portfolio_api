package service

import (
	"context"
	"fmt"

	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/store"
)

type DeleteExercise struct {
	DB   store.Execer
	Repo ExerciseDeleter
}

func (d *DeleteExercise) DeleteExercise(ctx context.Context, id entity.ExerciseID) error {
	affected, err := d.Repo.DeleteExercise(ctx, d.DB, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	if *affected != 1 {
		return fmt.Errorf("failed to delete: affected rows: %d", *affected)
	}
	return nil
}
