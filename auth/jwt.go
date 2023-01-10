package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/t-tazy/my_portfolio_api/clock"
	"github.com/t-tazy/my_portfolio_api/entity"
)

// JWTの独自クレーム名
const (
	RoleKey     = "role"
	UserNameKey = "user_name"
)

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

// 鍵として読み込んだデータを保持する
type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store                 Store
	Clocker               clock.Clocker // JWTの時刻情報を操作
}

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

// 鍵として読み込んだデータをJWTer型として保持
func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{Store: s}
	privkey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: private key: %w", err)
	}
	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: public key: %w", err)
	}
	j.PrivateKey = privkey
	j.PublicKey = pubkey
	j.Clocker = c
	return j, nil
}

// 単一のキーを解析
func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}
	return key, nil
}

// 引数で渡されたユーザーに対して署名済みのJWTを発行する
func (j *JWTer) GenerateToken(ctx context.Context, u *entity.User) ([]byte, error) {
	token, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/t-tazy/my_portfolio_api`).
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30*time.Minute)).
		Claim(RoleKey, u.Role).     // 独自クレーム
		Claim(UserNameKey, u.Name). // 独自クレーム
		Build()
	if err != nil {
		return nil, fmt.Errorf("GenerateToken: failed to build token: %w", err)
	}
	// Redisに保存
	if err := j.Store.Save(ctx, token.JwtID(), u.ID); err != nil {
		return nil, err
	}

	// 署名
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, err
	}
	return signed, nil
}

// HTTPリクエストからJWTを取得する
func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(
		r,
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, err
	}
	// 期限切れではないことを確認
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("GetToken: failed to validate token: %w", err)
	}
	// Redisから削除されていないか確認
	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("GetToken: %q expired: %w", token.JwtID(), err)
	}
	return token, nil
}

type userIDKey struct{} // contextのkeyとして使う独自型
type roleKey struct{}   // contextのkeyとして使う独自型

// contextにユーザーIDを付加
func SetUserID(ctx context.Context, uid entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

// contextからユーザーIDを取得
func GetUserID(ctx context.Context) (entity.UserID, bool) {
	id, ok := ctx.Value(userIDKey{}).(entity.UserID)
	return id, ok
}

// tokenからロール情報を取り出し、contextに付加
func SetRole(ctx context.Context, token jwt.Token) context.Context {
	get, ok := token.Get(RoleKey)
	if !ok {
		return context.WithValue(ctx, roleKey{}, "")
	}
	return context.WithValue(ctx, roleKey{}, get)
}

// contextからユーザーロールを取得
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey{}).(string)
	return role, ok
}

// リクエストスコープな値をcontextに含める
func (j *JWTer) FillContext(r *http.Request) (*http.Request, error) {
	token, err := j.GetToken(r.Context(), r)
	if err != nil {
		return nil, err
	}
	uid, err := j.Store.Load(r.Context(), token.JwtID())
	if err != nil {
		return nil, err
	}
	ctx := SetUserID(r.Context(), uid)
	ctx = SetRole(ctx, token)
	clone := r.Clone(ctx)
	return clone, nil
}

// contextから管理者権限の有無を確認する
func IsAdmin(ctx context.Context) bool {
	role, ok := GetRole(ctx)
	if !ok {
		return false
	}
	return role == "admin"
}
