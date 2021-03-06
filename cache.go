package kessaku

import (
	"math/rand"
)

type cache struct {
	freeWorkers []*worker
}

func NewCache(poolRef *Pool) *cache {
	workers := make([]*worker, 0)
	workers = append(workers, &worker{
		pool:  poolRef,
		tasks: make(chan func()),
	})
	return &cache{
		freeWorkers: workers,
	}
}

func (c *cache) Put(w *worker) {
	c.freeWorkers = append(c.freeWorkers, w)
}

func (c *cache) Get() *worker {
	if len(c.freeWorkers) == 0 {
		return nil
	}
	if len(c.freeWorkers) == 1 {
		w := c.freeWorkers[0]
		return w
	}
	ix := rand.Intn(len(c.freeWorkers) - 1)
	w := c.freeWorkers[ix]
	c.freeWorkers[len(c.freeWorkers)-1], c.freeWorkers[ix] = c.freeWorkers[ix], c.freeWorkers[len(c.freeWorkers)-1]
	c.freeWorkers = c.freeWorkers[:len(c.freeWorkers)-1]
	return w
}
