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

	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("processData error: %s\n", err)
			}
			fmt.Printf("processData finished at %s\n", time.Now())
			return

		case num := <-dataCh:
			fmt.Printf("processing data: %d\n", num)
		}
	}
}

func main() {
	ctx := context.Background()

	currentTime := time.Now()
	fmt.Printf("main started at %s\n", currentTime)

	deadline := currentTime.Add(2 * time.Second)
	ctxDeadline, cancelFunc := context.WithDeadline(ctx, deadline)
	defer cancelFunc()

	dataCh := make(chan int)
	go processData(ctxDeadline, dataCh)

	for i := 0; i <= 5; i++ {
		select {
		case dataCh <- i:
			time.Sleep(1 * time.Second)
		case <-ctxDeadline.Done():
			break
		}
	}

	// cancelFunc()

	select {
	case <-ctxDeadline.Done():
		fmt.Println("main context canceled")
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("main finished at %s\n", time.Now())
}
