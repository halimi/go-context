package main

import (
	"context"
	"fmt"
)

func doAnotherThing(ctx context.Context) {
	fmt.Printf("doAnotherThing: myKey's value: %s\n", ctx.Value("myKey"))
}

func doSomething(ctx context.Context) {
	fmt.Printf("doSomething: myKey's value: %s\n", ctx.Value("myKey"))

	ctxAnother := context.WithValue(ctx, "myKey", "anotherValue")
	doAnotherThing(ctxAnother)

	fmt.Printf("doSomething: myKey's value: %s\n", ctx.Value("myKey"))
}

func main() {
	ctx := context.Background()

	ctxValue := context.WithValue(ctx, "myKey", "myValue")

	doSomething(ctxValue)
}
