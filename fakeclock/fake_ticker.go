package fakeclock

import (
	"sync"
	"time"

	"github.com/pivotal-golang/clock"
)

type fakeTicker struct {
	clock clock.Clock

	mutex    sync.Mutex
	duration time.Duration
	channel  chan time.Time

	timer clock.Timer
}

func NewFakeTicker(clock clock.Clock, d time.Duration) clock.Ticker {
	channel := make(chan time.Time)
	timer := clock.NewTimer(d)

	ft := fakeTicker{
		clock:    clock,
		duration: d,
		channel:  channel,
		timer:    timer,
	}

	go ft.run()

	return &ft
}

func (ft *fakeTicker) run() {
	for {
		time := <-ft.timer.C()
		ft.timer = ft.clock.NewTimer(ft.duration)
		ft.channel <- time
	}
}

func (ft *fakeTicker) C() <-chan time.Time {
	ft.mutex.Lock()
	defer ft.mutex.Unlock()
	return ft.channel
}

func (ft *fakeTicker) Stop() {
	ft.mutex.Lock()
	ft.timer.Stop()
	ft.mutex.Unlock()
}
