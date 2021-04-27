package main

import (
	"fmt"
	"time"

	k "github.com/amoghyermalkar123/kessaku"
)

func main() {
	pool, err := k.NewPool(k.WithContext(), k.WithPoolSize(2))

	if err != nil {
		return
	}

	pool.Submit(func() {
		fmt.Println("ONE TASK")
	})

	pool.Submit(func() {
		fmt.Println("TWO TASK")
	})

	pool.Submit(func() {
		fmt.Println("THREE TASK")
	})

	pool.Submit(func() {
		fmt.Println("FOUR TASK")
	})
	pool.Submit(func() {
		fmt.Println("FIVE TASK")
	})
	pool.Submit(func() {
		fmt.Println("SIX TASK")
	})
	rw := pool.AtCapacity
	fmt.Println(rw)

	time.Sleep(time.Second * 6)

	rw = pool.AtCapacity
	fmt.Println(rw)

}
