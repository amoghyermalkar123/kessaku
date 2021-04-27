package kessaku_test

import (
	"log"
	"time"

	k "github.com/amoghyermalkar123/kessaku"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pool Tests", func() {

	Describe("NewPool", func() {

		It("Should create a pool with arbitrary amount of workers", func() {
			pool, err := k.NewPool(k.WithPoolSize(5))
			Expect(err).To(BeNil())
			Expect(pool.Opts.PoolSize).To(Equal(5))
			Expect(pool.Opts.WithContext).To(Equal(false))
		})

		It("Should be able to submit tasks to the pool", func() {
			pool, err := k.NewPool(k.WithPoolSize(5))
			Expect(err).To(BeNil())
			Expect(pool.Opts.PoolSize).To(Equal(5))
			Expect(pool.Opts.WithContext).To(Equal(false))

			err = pool.Submit(func() {
				log.Println("Im a task!!")
			})
			Expect(err).To(BeNil())

			err = pool.Submit(func() {
				log.Println("Im a task2!!")
			})
			Expect(err).To(BeNil())

			err = pool.Submit(func() {
				log.Println("Im a task3!!")
			})
			Expect(err).To(BeNil())
			time.Sleep(time.Second)
		})

	})
})
