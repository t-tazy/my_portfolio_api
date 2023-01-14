package handler

import (
	"context"

	"github.com/t-tazy/my_portfolio_api/entity"
)

// リクエストの解釈とレスポンスの構築以外を次のインターフェースに移譲

//go:generate go run github.com/matryer/moq -out moq_test.go . ListExercisesService AddExerciseService GetExerciseService DeleteExerciseService UpdateExerciseService RegisterUserService LoginService
type ListExercisesService interface {
	ListExercises(ctx context.Context) (entity.Exercises, error)
}

type AddExerciseService interface {
	AddExercise(ctx context.Context, title, description string) (*entity.Exercise, error)
}

type GetExerciseService interface {
	GetExercise(ctx context.Context, id entity.ExerciseID) (*entity.Exercise, error)
}

type DeleteExerciseService interface {
	DeleteExercise(ctx context.Context, id entity.ExerciseID) error
}

type UpdateExerciseService interface {
	UpdateExercise(ctx context.Context, id entity.ExerciseID, title, description string) error
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}
