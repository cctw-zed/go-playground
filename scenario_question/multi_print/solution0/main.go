package main

import (
	"fmt"
	"sync"
)

func main() {

	chOdd := make(chan struct{})
	chEven := make(chan struct{})

	ma := make(map[int]struct{})
	chMap := make(chan int, 100)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			<-chEven
			chMap <- i * 2
			chOdd <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			<-chOdd
			chMap <- i*2 + 1
			chEven <- struct{}{}
		}
	}()

	go func() {
		// 启动
		chEven <- struct{}{}

		wg.Wait()
		close(chMap)
	}()

	for v := range chMap {
		ma[v] = struct{}{}
		fmt.Println(v)
	}

	fmt.Printf("map: %v", ma)
}
