package common

import "time"

func Time[T any](runnable func() (T, error)) (time.Duration, T, error) {
	t := time.Now()
	res, err := runnable()
	return time.Since(t), res, err
}
