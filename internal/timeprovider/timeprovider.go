package timeprovider

import "time"

type TimeProvider interface {
	Now() time.Time
}

type timeFunc func() time.Time

func (t timeFunc) Now() time.Time {
	return t()
}

func NewSystemTimeProvider() TimeProvider {
	return timeFunc(time.Now)
}
