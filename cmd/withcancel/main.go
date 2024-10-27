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
			fmt.Println("processData finished")
			return

		case num := <-dataCh:
			fmt.Printf("processing data: %d\n", num)
		}
	}
}

func main() {
	ctx := context.Background()

	ctxCancel, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	dataCh := make(chan int)
	go processData(ctxCancel, dataCh)

	for i := 0; i <= 5; i++ {
		dataCh <- i
	}

	cancelFunc()

	select {
	case <-ctxCancel.Done():
		fmt.Println("main context canceled")
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("main finished")
}
