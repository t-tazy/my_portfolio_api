package store

import (
	"context"

	"github.com/t-tazy/my_portfolio_api/entity"
)

// 引数で受け取ったユーザーIDのエクササイズを取得する
func (r *Repository) ListExercises(
	ctx context.Context, db Queryer, id entity.UserID,
) (entity.Exercises, error) {
	exercises := entity.Exercises{}
	sql := `SELECT
			id, user_id, title, description, created, modified
			FROM exercises
			WHERE user_id = ?;`
	if err := db.SelectContext(ctx, &exercises, sql, id); err != nil {
		return nil, err
	}
	return exercises, nil
}

// エクササイズを保存する
// 引数として受け取った*entity.Exercise型のIDフィールドを更新することで
// 呼び出し元にRDBMSより発行されたIDを伝える
func (r *Repository) AddExercise(
	ctx context.Context, db Execer, e *entity.Exercise,
) error {
	e.Created = r.Clocker.Now()
	e.Modified = r.Clocker.Now()
	sql := `INSERT INTO exercises
			(user_id, title, description, created, modified)
			VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(
		ctx, sql, e.UserID, e.Title, e.Description, e.Created, e.Modified,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = entity.ExerciseID(id)
	return nil
}

// エクササイズを取得する
func (r *Repository) GetExercise(
	ctx context.Context, db Queryer, id entity.ExerciseID,
) (*entity.Exercise, error) {
	exercise := &entity.Exercise{}
	sql := `SELECT * FROM exercises WHERE id = ?;`
	if err := db.GetContext(ctx, exercise, sql, id); err != nil {
		return nil, err
	}
	return exercise, nil
}

// エクササイズを削除する
func (r *Repository) DeleteExercise(
	ctx context.Context, db Execer, id entity.ExerciseID,
) (*int64, error) {
	sql := `DELETE FROM exercises WHERE id = ?;`
	result, err := db.ExecContext(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &affected, nil
}

// エクササイズを更新する
func (r *Repository) UpdateExercise(
	ctx context.Context, db Execer, e *entity.Exercise,
) (*int64, error) {
	e.Modified = r.Clocker.Now()
	sql := `UPDATE exercises SET
			title = ?, description = ?, modified = ?
			WHERE id = ?;`
	result, err := db.ExecContext(
		ctx, sql, e.Title, e.Description, e.Modified, e.ID,
	)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &affected, nil
}
