package main

import (
	"context"
	"fmt"
	"time"
)

func x111() *time.Timer {
	// timer到达后执行，某个函数；timer提前取消则不会执行
	return time.AfterFunc(3*time.Second, func() {
		fmt.Println("xxxxxxx")
	})
}

func WithSignal(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {

}

func x112() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tctx, tcel := context.WithTimeout(ctx, 6*time.Second)
	defer tcel()

	context.WithCancel(tctx)
	// context.WithDeadline(tctx, time.Now())
	// context.WithValue(tctx, "key", "value")
}

func main() {
	//
	tm := x111()
	tm.Stop()

	time.Sleep(6 * time.Second)

}
