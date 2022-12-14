package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/testutil"
)

func TestAddExercise(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_exercise/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_exercise/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_exercise/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_exercise/bad_rsp.json.golden",
			},
		},
	}
	for key, test := range tests {
		// クロージャ用に変数を束縛
		test := test
		t.Run(key, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/exercises",
				bytes.NewReader(testutil.LoadFile(t, test.reqFile)),
			)

			moq := &AddExerciseServiceMock{}
			moq.AddExerciseFunc = func(
				ctx context.Context, title, description string,
			) (*entity.Exercise, error) {
				if test.want.status == http.StatusOK {
					return &entity.Exercise{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := AddExercise{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)
			rsp := w.Result()
			testutil.AssertResponse(
				t, rsp, test.want.status, testutil.LoadFile(t, test.want.rspFile),
			)
		})
	}
}
