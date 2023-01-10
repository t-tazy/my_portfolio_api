package clock

import "time"

type Clocker interface {
	Now() time.Time
}

// 現在時刻
type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

// テスト用の固定時刻
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 12, 22, 9, 34, 0, 0, time.UTC)
}
