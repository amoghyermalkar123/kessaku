package kessaku

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Batch struct {
	stopper   func()
	processor *batchWorker
}

type batchWorker struct {
	tasks      chan func()
	isInactive bool
}

func (b *batchWorker) batch(ctx context.Context) {
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
