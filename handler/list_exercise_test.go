package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/testutil"
)

func TestListExercise(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		exercises []*entity.Exercise
		want      want
	}{
		"ok": {
			exercises: []*entity.Exercise{
				{
					ID:          1,
					Title:       "test1",
					Description: "",
				},
				{
					ID:          2,
					Title:       "test2",
					Description: "テスト用",
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_exercise/ok_rsp.json.golden",
			},
		},
		"empty": {
			exercises: []*entity.Exercise{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_exercise/empty_rsp.json.golden",
			},
		},
	}
	for key, test := range tests {
		test := test
		t.Run(key, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/exercises", nil)

			moq := &ListExercisesServiceMock{}
			moq.ListExercisesFunc = func(
				ctx context.Context,
			) (entity.Exercises, error) {
				if test.exercises != nil {
					return test.exercises, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := ListExercise{Service: moq}
			sut.ServeHTTP(w, r)
			rsp := w.Result()
			testutil.AssertResponse(
				t, rsp, test.want.status, testutil.LoadFile(t, test.want.rspFile),
			)
		})
	}
}
