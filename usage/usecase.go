package main

import (
	"fmt"
	"time"

	k "github.com/amoghyermalkar123/kessaku"
)

func main() {
	pool, err := k.NewPool(k.WithContext(false), k.WithPoolSize(2))

	if err != nil {
		return
	}

	pool.Submit(func() {
		fmt.Println("ONE TASK")
	})

	pool.Submit(func() {
		fmt.Println("TWO TASK")
	})

	time.Sleep(time.Second * 4)
	pool.Submit(func() {
		fmt.Println("THREE TASK")
	})
}
