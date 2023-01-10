package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserID int

type User struct {
	ID       UserID    `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Password string    `json:"password" db:"password"`
	Role     string    `json:"role" db:"role"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

// ハッシュ化されて永続化されていたパスワードと入力されたパスワードを比較検証する
func (u *User) ComparePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
}
