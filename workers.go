package kessaku

import "log"

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
	log.Println("worker started")
	go func() {

		for fn := range w.tasks {
			fn()
		}

		w.pool.cache.Put(w)
	}()
	// defer recovery logic
	// run the tasks
	// after the tasks are done
	// put the worker back in the pool
	// put the worker in the cache

}
