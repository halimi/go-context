package main

import (
	"context"
	"fmt"
	"time"
)

func processData(ctx context.Context, dataCh chan int) {
	deadline, ok := ctx.Deadline()
	if ok {
		fmt.Printf("we have a deadline: %s\n", deadline)
	} else {
		fmt.Println("no deadline!")
	}

	for num := 0; num <= 5; num++ {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("processData error: %s\n", err)
			}
			fmt.Printf("processData finished at %s\n", time.Now())
			return

		case dataCh <- num:
			fmt.Printf("processing data: %d\n", num)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx := context.Background()

	currentTime := time.Now()
	fmt.Printf("main started at %s\n", currentTime)

	ctxTimeout, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	dataCh := make(chan int)
	go processData(ctxTimeout, dataCh)

	stopped := false
	for !stopped {
		select {
		case result := <-dataCh:
			fmt.Printf("received result: %d\n", result)
		case <-ctxTimeout.Done():
			fmt.Println("main context canceled")
			stopped = true
		}
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("main finished at %s\n", time.Now())
}
