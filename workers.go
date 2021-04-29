package kessaku

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type worker struct {
	pool  *Pool
	tasks chan func()
}

func NewWorker(p *Pool) *worker {
	return &worker{
		pool:  p,
		tasks: make(chan func()),
	}
}

func (w *worker) run() {
	go func() {
		// panic recovery
		defer func() {
			if p := recover(); p != nil {
				log.Warn("worker exiting with a panic %s", p)
			}

		}()

		for {
			select {
			case <-time.After(5 * time.Second):
				w.pool.cache.Put(w)
				w.pool.AtCapacity--
				return
			case fn := <-w.tasks:
				if fn == nil {
					return
				}
				fn()
			}
		}
	}()
}
