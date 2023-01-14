package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/t-tazy/my_portfolio_api/entity"
)

type GetExercise struct {
	Service GetExerciseService
}

func (ge *GetExercise) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// パスパラメータ取得
	exerciseID, err := strconv.Atoi(chi.URLParam(r, "exerciseID"))
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	eid := entity.ExerciseID(exerciseID)

	exercise, err := ge.Service.GetExercise(ctx, eid)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID          entity.ExerciseID `json:"id"`
		UserID      entity.UserID     `json:"user_id"`
		Title       string            `json:"title"`
		Description string            `json:"description"`
	}{exercise.ID, exercise.UserID, exercise.Title, exercise.Description}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
