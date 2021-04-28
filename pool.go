package kessaku

import (
	"context"
	"errors"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Pool struct {
	AtCapacity   int
	Opts         *Options
	cache        *cache
	batchManager *Batch
}

var initBatchWorker sync.Once

func NewPool(options ...OptionSetter) (*Pool, error) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	// options pattern
	opts := loadOptions(options...)
	// pool instance
	pool := &Pool{Opts: opts}
	// the cache
	pool.cache = NewCache(pool)
	return pool, nil
}

func (p *Pool) Submit(task func()) error {
	// process order respecting bactched jobs
	if p.Opts.WithBatch && !p.batchManager.processor.isInactive {
		p.batchProcessor(task)
	}
	// get from cache
	w := p.cache.Get()
	// use the worker thats free
	if w != nil {
		w.run()
		w.tasks <- task
		p.AtCapacity++
		return nil
	}
	// if no worker is free, check if pool has reached size
	if p.AtCapacity == p.Opts.PoolSize {
		log.Info("Pool has reached capacity, remember you cant create workers for a while")
		return errors.New("no worker is free and pool has reached capacity")
	}
	// if pool is not at capacity create one at run-ime
	w = NewWorker(p)
	// spawn the newly created worker
	p.spawnWorkerAndUpdateCapacity(w, task)
	return nil
}

func (p *Pool) spawnWorkerAndUpdateCapacity(w *worker, task func()) {
	w.run()
	w.tasks <- task
	p.AtCapacity++
}

func (p *Pool) RunningWorkers() int {
	return p.AtCapacity
}

func (p *Pool) batchProcessor(task func()) {
	if p.batchManager.stopper == nil {
		ctx := context.Background()
		ctx, stopBatchedWorker := context.WithCancel(ctx)
		p.batchManager.stopper = stopBatchedWorker
		p.batchManager.processor = &batchWorker{}
		p.batchManager.processor.tasks = make(chan func())
		go p.batchManager.processor.batch(ctx)
		p.batchManager.processor.tasks <- task
	} else {
		p.batchManager.processor.tasks <- task
	}
}
