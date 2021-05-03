package kessaku

import (
	"context"
	"os"

	e "github.com/amoghyermalkar123/kessaku/errors"
	log "github.com/sirupsen/logrus"
)

// Pool type
type Pool struct {
	AtCapacity int
	Opts       *Options
	cache      *cache
	batch      *batchWorker
}

type runner interface {
	run(context.Context)
}

func runProcs(r runner, ctx context.Context) {
	r.run(ctx)
}

// NewPool instantiates new Pool instance
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

// Submit proides an API to get tasks as an input parameter and submit those to concurrent workers or batch workers
// based on the options set by the user
func (p *Pool) Submit(task func()) error {
	// ctx
	ctx := context.Background()
	// process order respecting bactched jobs
	if p.Opts.WithBatch {
		p.batchProcessor(task)
		return nil
	}
	// get from cache
	w := p.cache.Get()
	// use the worker thats free
	if w != nil {
		runProcs(w, ctx)
		w.tasks <- task
		p.AtCapacity++
		return nil
	}
	// if no worker is free, check if pool has reached size
	if p.AtCapacity == p.Opts.PoolSize {
		log.Info("Pool has reached capacity, remember you cant create workers for a while")
		return e.ErrPoolReachedCapacity
	}
	// if pool is not at capacity create one at run-ime
	w = NewWorker(p)
	// spawn the newly created worker
	p.spawnWorkerAndUpdateCapacity(w, task, ctx)
	return nil
}

func (p *Pool) spawnWorkerAndUpdateCapacity(w *worker, task func(), ctx context.Context) {
	runProcs(w, ctx)
	w.tasks <- task
	p.AtCapacity++
}

// RunningWorkers returns number of running concurrent workers
func (p *Pool) RunningWorkers() int {
	return p.AtCapacity
}

func (p *Pool) batchProcessor(task func()) {
	if b, ctx := NewBatchWorker(); b.stopper != nil {
		// new instance creation either for the first time or after a stop
		p.batch = b
		go runProcs(b, ctx)
		b.tasks <- task
	} else {
		// use the already created batch worker ref stored in pool
		p.batch.tasks <- task
	}
}

func (p *Pool) StopBatchWorker() error {
	p.batch.stopper()
	p.batch.isInactive = true
	p.batch = nil
	return nil
}
