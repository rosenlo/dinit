package worker

import (
	"context"
	"time"
)

type Worker interface {
	Start(run func())
	Shutdown()
}

type Runner struct {
	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRunner(ctx context.Context, interval int) (client Worker, err error) {
	if err != nil {
		return
	}
	client = &Runner{
		ticker: time.NewTicker(time.Duration(interval) * time.Second),
		ctx:    ctx,
	}
	return
}

func (s *Runner) Start(run func()) {
	for {
		select {
		case <-s.ticker.C:
			run()
		case <-s.ctx.Done():
			s.Shutdown()
			return
		}
	}
}

func (s *Runner) Shutdown() {
	s.ticker.Stop()
}
