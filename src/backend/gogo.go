package gogo

import (
	"sync"
	"time"
)

// Runs fn at the specified time.
func At(t time.Time, fn func()) (done <-chan bool) {
	return After(t.Sub(time.Now()), fn)
}

// Runs until time in every dur.
func Until(t time.Time, dur time.Duration, fn func()) (done <-chan bool) {
	ch := make(chan bool, 1)
	untilRecv(ch, t, dur, fn)
	return ch
}

func untilRecv(ch chan bool, t time.Time, dur time.Duration, fn func()) {
	if t.Sub(time.Now()) > 0 {
		time.AfterFunc(dur, func() {
			fn()
			untilRecv(ch, t, dur, fn)
		})
		return
	}
	ch <- true
}

// Runs fn after duration. Similar to time.AfterFunc
func After(duration time.Duration, fn func()) (done <-chan bool) {
	ch := make(chan bool, 1)
	time.AfterFunc(duration, func() {
		fn()
		ch <- true
	})
	return ch
}

func main() {

}
