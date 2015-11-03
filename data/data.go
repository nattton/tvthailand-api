package data

const (
	limitRow       = 40
	maxConcurrency = 8
)

var throttle = make(chan int, maxConcurrency)
