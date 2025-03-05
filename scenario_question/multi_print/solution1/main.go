package main

import (
	"fmt"
	"sync"
)

func main() {
	maxNumber := 100
	var wg sync.WaitGroup

	// 创建两个通道用于同步
	oddTurn := make(chan bool)
	evenTurn := make(chan bool)

	res := make(map[int]bool)
	mu := sync.Mutex{}

	wg.Add(2)

	// 打印偶数的协程
	go func() {
		defer wg.Done()

		for i := 0; i <= maxNumber; i += 2 {
			// 等待轮到偶数打印
			<-evenTurn
			mu.Lock()
			res[i] = true
			mu.Unlock()

			// 通知奇数协程可以打印了（仅当尚未达到最大值时）
			if i+1 <= maxNumber {
				oddTurn <- true
			}
		}
	}()

	// 打印奇数的协程
	go func() {
		defer wg.Done()

		for i := 1; i <= maxNumber; i += 2 {
			// 等待轮到奇数打印
			<-oddTurn
			mu.Lock()
			res[i] = true
			mu.Unlock()

			// 通知偶数协程可以打印了（仅当尚未达到最大值时）
			if i+1 <= maxNumber {
				evenTurn <- true
			}
		}
	}()

	// 启动交替打印过程，从偶数开始
	evenTurn <- true

	// 等待两个协程完成
	wg.Wait()

	fmt.Println(res)

	fmt.Println("所有数字都已打印完毕")
}
