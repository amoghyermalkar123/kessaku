package kessaku

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

type Pool struct {
	AtCapacity int
	Opts       *Options
	cache      *cache
	// TODO: implement TaskCache
}

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
