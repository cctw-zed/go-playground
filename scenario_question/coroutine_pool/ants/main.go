package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"time"
)

func main() {
	p, _ := ants.NewPool(100)
	defer p.Release()

	for i := 0; i < 1000; i++ {
		err := p.Submit(func() {
			print(i)
		})
		if err != nil {
			panic(err)
		}
	}
}

func print(x int) {
	fmt.Println(x)
	time.Sleep(time.Second * 10)
}
