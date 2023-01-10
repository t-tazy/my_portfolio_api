package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/t-tazy/my_portfolio_api/clock"
	"github.com/t-tazy/my_portfolio_api/entity"
	"github.com/t-tazy/my_portfolio_api/store"
	"github.com/t-tazy/my_portfolio_api/testutil/fixture"
)

// 鍵ファイルが埋め込めているかテスト
func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}

	want = []byte("-----BEGIN PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s, but got %s", want, rawPrivKey)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	ctx := context.Background()
	moq := &StoreMock{}
	wantID := entity.UserID(10)
	u := fixture.User(&entity.User{ID: wantID}) // ダミー生成
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID, userID)
		}
		return nil
	}
	sut, err := NewJWTer(moq, clock.FixedClocker{})
	if err != nil {
		t.Fatal(err)
	}
	got, err := sut.GenerateToken(ctx, u)
	if err != nil {
		t.Fatalf("not want err: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("token is empty")
	}
}

// 正常系
func TestJWTer_GetToken(t *testing.T) {
	t.Parallel()

	c := clock.FixedClocker{}
	want, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/t-tazy/my_portfolio_api`).
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").          // 独自クレーム
		Claim(UserNameKey, "test_user"). // 独自クレーム
		Build()
	if err != nil {
		t.Fatal(err)
	}
	// 鍵を解析
	privkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}
	// 署名
	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, privkey))
	if err != nil {
		t.Fatal(err)
	}

	userID := entity.UserID(10) // ダミーデータ
	ctx := context.Background()
	moq := &StoreMock{}
	moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
		return userID, nil
	}
	sut, err := NewJWTer(moq, c)
	if err != nil {
		t.Fatal(err)
	}

	// HTTPリクエストの作成
	req := httptest.NewRequest(http.MethodGet, "https://github.com/t-tazy", nil)
	// HTTPリクエストヘッダーにアクセストークンを付与
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signed))

	got, err := sut.GetToken(ctx, req)
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetToken: got = %v, want %v", got, want)
	}
}

// テスト用の固定時刻+24時間
type FixedTomorrowClocker struct{}

func (c FixedTomorrowClocker) Now() time.Time {
	return clock.FixedClocker{}.Now().Add(24 * time.Hour)
}

// 異常系
func TestJWTer_GetToken_NG(t *testing.T) {
	t.Parallel()

	c := clock.FixedClocker{}
	token, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/t-tazy/my_portfolio_api`).
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").          // 独自クレーム
		Claim(UserNameKey, "test_user"). // 独自クレーム
		Build()
	if err != nil {
		t.Fatal(err)
	}
	// 鍵を解析
	privkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}
	// 署名
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privkey))
	if err != nil {
		t.Fatal(err)
	}

	type moq struct {
		userID entity.UserID
		err    error
	}
	tests := map[string]struct {
		c   clock.Clocker
		moq moq
	}{
		// トークン期限切れ
		"expire": {
			c: FixedTomorrowClocker{},
		},
		// Redis上にデータが存在しない
		"notFoundInStore": {
			c: clock.FixedClocker{},
			moq: moq{
				err: store.ErrNotFound,
			},
		},
	}
	for key, test := range tests {
		// クロージャ用に変数を束縛
		test := test
		t.Run(key, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			moq := &StoreMock{}
			moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
				return test.moq.userID, test.moq.err
			}
			sut, err := NewJWTer(moq, test.c)
			if err != nil {
				t.Fatal(err)
			}

			// HTTPリクエストの作成
			req := httptest.NewRequest(http.MethodGet, "https://github.com/t-tazy", nil)
			// HTTPリクエストヘッダーにアクセストークンを付与
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signed))

			got, err := sut.GetToken(ctx, req)
			if err == nil {
				t.Error("want error, but got nil")
			}
			if got != nil {
				t.Errorf("want nil, but got %v", got)
			}
		})
	}
}
