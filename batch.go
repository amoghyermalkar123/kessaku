package kessaku

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type batchWorker struct {
	stopper func()
	tasks   chan func()
	// currently there's not much use of this variable but will be keeping it in-case i ever need it
	isInactive bool
}

var once sync.Once

func NewBatchWorker() (*batchWorker, context.Context) {
	bw := &batchWorker{}
	ctx := context.Background()
	ctx, stopBatchedWorker := context.WithCancel(ctx)

	once.Do(func() {
		bw.stopper = stopBatchedWorker
		bw.tasks = make(chan func())
	})

	return bw, ctx
}

func (b *batchWorker) run(ctx context.Context) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				logrus.Warn("The batch processor has panicked, your further task submissions will now be delivered to concurrent workers")
				b.isInactive = true
			}
		}()

		for {
			select {
			case <-ctx.Done():
				logrus.Info("Batched Jobs have stopped")
				return
			case task := <-b.tasks:
				task()
			}
		}
	}()
}
