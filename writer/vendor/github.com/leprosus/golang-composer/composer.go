package golang_composer

import (
	"sync"
	"sync/atomic"
)

type composer struct {
	flag uint64
	chn  chan bool
}

var (
	instance *composer
	once     sync.Once
)

func GetComposer() *composer {
	once.Do(func() {
		chn := make(chan bool)
		close(chn)

		instance = &composer{chn: chn}
	})

	return instance
}

// Checks rather to wait or does not need
func (c *composer) NeedWait() {
	<-c.chn
}

// Says all of goroutins to resume execution
func (c *composer) Play() {
	if atomic.LoadUint64(&c.flag) == 1 {
		close(c.chn)
		atomic.StoreUint64(&c.flag, 0)
	}
}

// Says all of goroutins to pause execution
func (c *composer) Pause() {
	if atomic.LoadUint64(&c.flag) == 0 {
		atomic.StoreUint64(&c.flag, 1)
		c.chn = make(chan bool)
	}
}
