package handler

import (
	"context"

	"github.com/t-tazy/my_portfolio_api/entity"
)

// リクエストの解釈とレスポンスの構築以外を次のインターフェースに移譲

//go:generate go run github.com/matryer/moq -out moq_test.go . ListExercisesService AddExerciseService RegisterUserService LoginService
type ListExercisesService interface {
	ListExercises(ctx context.Context) (entity.Exercises, error)
}

type AddExerciseService interface {
	AddExercise(ctx context.Context, title, description string) (*entity.Exercise, error)
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}
