package concurrent

import "sync"

type LimitedWaitGroup struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

func NewLimitedWaitGroup(maxConcurrency int) *LimitedWaitGroup {
	return &LimitedWaitGroup{
		sem: make(chan struct{}, maxConcurrency),
	}
}

func (l *LimitedWaitGroup) Go(fn func()) {
	l.sem <- struct{}{}
	l.wg.Go(func() {
		defer func() { <-l.sem }()

		fn()
	})
}

func (l *LimitedWaitGroup) Wait() {
	l.wg.Wait()
}
