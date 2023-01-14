package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/t-tazy/my_portfolio_api/auth"
	"github.com/t-tazy/my_portfolio_api/entity"
)

type DeleteExercise struct {
	Service DeleteExerciseService
	*GetExercise
}

func (de *DeleteExercise) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// エクササイズを取得
	exercise, err := de.GetExercise.Service.GetExercise(ctx, eid)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	id, ok := auth.GetUserID(ctx)
	if !ok || id != exercise.UserID {
		RespondJSON(ctx, w, &ErrResponse{
			Message: fmt.Errorf("failed to delete: permission denied").Error(),
		}, http.StatusBadRequest)
		return
	}

	if err := de.Service.DeleteExercise(ctx, eid); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID entity.ExerciseID `json:"id"`
	}{ID: eid}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
