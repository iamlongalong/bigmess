package main

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

func RetryTimes(f func() error, maxTimes int, baseSpan time.Duration) error {
	t := baseSpan
	var err error

	for i := 0; i < maxTimes; i++ {
		err = f()
		if err != nil {
			time.Sleep(t * time.Duration(i))
			continue
		}

		return nil
	}

	return err
}

func UUIDV4() string {
	return uuid.NewV4().String()
}
