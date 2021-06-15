package main

import (
	"context"
	"fmt"
	"time"
)

type signalCtx struct {
	context.Context
}

/*
Deadline() (deadline time.Time, ok bool)
Done() <-chan struct{}
Err() error
Value(key interface{}) interface{}
*/

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}

func main() {
	ctx := context.Background()
	tctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	go handle(tctx, 3*time.Microsecond)

	select {
	case <-tctx.Done():
		fmt.Println("handle", tctx.Err())
	}

	context.WithCancel(ctx)

}
