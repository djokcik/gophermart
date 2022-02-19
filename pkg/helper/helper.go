package helpers

import (
	"time"
)

func SetTicker(fn func(), seconds time.Duration) {
	ticker := time.NewTicker(seconds)

	for range ticker.C {
		fn()
	}
}
