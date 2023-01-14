package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/t-tazy/my_portfolio_api/auth"
	"github.com/t-tazy/my_portfolio_api/entity"
)

type UpdateExercise struct {
	Service   UpdateExerciseService
	Validator *validator.Validate
	*GetExercise
}

func (ue *UpdateExercise) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	exercise, err := ue.GetExercise.Service.GetExercise(ctx, eid)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	id, ok := auth.GetUserID(ctx)
	if !ok || id != exercise.UserID {
		RespondJSON(ctx, w, &ErrResponse{
			Message: fmt.Errorf("failed to update: permission denied").Error(),
		}, http.StatusBadRequest)
		return
	}

	var body struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// バリデーション
	if err := ue.Validator.Struct(body); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if err := ue.Service.UpdateExercise(ctx, eid, body.Title, body.Description); err != nil {
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
