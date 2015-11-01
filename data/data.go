package data

const (
	limitRow       = 40
	maxConcurrency = 4
)

var throttle = make(chan int, maxConcurrency)
