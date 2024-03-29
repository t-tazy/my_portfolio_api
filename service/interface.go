package service

import (
	"context"

	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/store"
)

// storeパッケージの直接参照を避ける

//go:generate go run github.com/matryer/moq -out moq_test.go . ExerciseAdder ExerciseLister ExerciseGetter ExerciseDeleter ExerciseUpdater UserRegister UserGetter TokenGenerator
type ExerciseAdder interface {
	AddExercise(ctx context.Context, db store.Execer, e *entity.Exercise) error
}

type ExerciseLister interface {
	ListExercises(ctx context.Context, db store.Queryer, id entity.UserID) (entity.Exercises, error)
}

type ExerciseGetter interface {
	GetExercise(ctx context.Context, db store.Queryer, id entity.ExerciseID) (*entity.Exercise, error)
}

type ExerciseDeleter interface {
	DeleteExercise(ctx context.Context, db store.Execer, id entity.ExerciseID) (*int64, error)
}

type ExerciseUpdater interface {
	UpdateExercise(ctx context.Context, db store.Execer, e *entity.Exercise) (*int64, error)
}

type UserRegister interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, name string) (*entity.User, error)
}

// *auth.JWTerの直接参照を避ける
type TokenGenerator interface {
	GenerateToken(ctx context.Context, u *entity.User) ([]byte, error)
}
