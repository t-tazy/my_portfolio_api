package entity

import "time"

type ExerciseID int

// 誤代入を防ぐため、IDフィールドは独自型を使う
type Exercise struct {
	ID          ExerciseID `json:"id" db:"id"`
	UserID      UserID     `json:"user_id" db:"user_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Created     time.Time  `json:"created" db:"created"`
	Modified    time.Time  `json:"modified" db:"modified"`
}

type Exercises []*Exercise
