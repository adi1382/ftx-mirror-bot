package main

import (
	"fmt"
	"github.com/adi1382/ftx-mirror-bot/client"
	"sync"
	"time"
)

var (
	subRoutineCloser chan int
	wg               sync.WaitGroup
)

func init() {
	subRoutineCloser = make(chan int, 100)
}

func main() {
	hostClient := client.NewHostClient("kqAyKxRHgQreYe4iNLB7qnpSp1zQsjQP2ePFUDjq", "PhqPf5qpoCp7aFjYC4Ua5ZJTAHuBP20P0TwyZvOX", subRoutineCloser, &wg)
	hostClient.Initialize()

	go func() {
		for {
			fmt.Println(hostClient.FetchOpenOrders())
			fmt.Println(hostClient.FetchOpenPositions())
			time.Sleep(time.Second)
		}
	}()

	//go func() {
	//	time.Sleep(time.Minute)
	//	fmt.Println("Attemptingggg")
	//	subRoutineCloser <- 0
	//}()

	wg.Wait()
	fmt.Println("wait group completed")

	//n := 0

	//go func() {
	//	for {
	//		fmt.Printf("\n\nActive Positions: %v\n", hostClient.ActivePositions())
	//		time.Sleep(time.Second * 5)
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		fmt.Printf("\n\nActive Orders: %v\n", hostClient.ActiveOrders())
	//		time.Sleep(time.Second * 5)
	//	}
	//}()

	//fmt.Println("$$$$$$$444")

	select {}
}
