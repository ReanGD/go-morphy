package benchmarks

import (
	"runtime"
	"time"
)

type benchmark struct {
	repeats int
	runs    int
	signal  chan interface{}
	f       func(int)
	result  float64
}

func (b *benchmark) launch() {
	defer func() {
		b.signal <- b
	}()

	first := true
	var minDuration time.Duration
	for i := 0; i < b.runs; i++ {
		runtime.GC()
		start := time.Now()
		b.f(b.repeats)
		duration := time.Now().Sub(start)
		if first {
			minDuration = duration
		} else {
			if minDuration > duration {
				minDuration = duration
			}
		}
	}
	b.result = minDuration.Seconds()
}

func runBenchmark(runs, repeats int, f func(int)) float64 {
	b := benchmark{
		repeats: repeats,
		runs:    runs,
		signal:  make(chan interface{}),
		f:       f,
	}

	go b.launch()
	<-b.signal
	return b.result
}
