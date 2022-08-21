package main

import (
	"context"
	"fmt"
	"github.com/czc-beijing/goredis/lib/sync/wait"
	"runtime"
	"time"
)

func dosomething(ctx context.Context, ch chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\ntimeout")
			return
		case out := <-ch:
			time.Sleep(5 * time.Second)
			fmt.Printf("%d\n", out)
		}
	}
}

func main() {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	lib := wait.Wait{}
	ctx, cancelFunc := context.WithCancel(context.Background())
	for i := 0; i < 5; i++ {
		lib.Add(1)
		go func(lib *wait.Wait, ch chan int) {
			dosomething(ctx, ch)
			defer lib.Done()
		}(&lib, ch)
	}
	_ = lib.WaitWithTimeout(4 * time.Second)
	cancelFunc()
	time.Sleep(4 * time.Second)
	fmt.Println("end")
	fmt.Println(runtime.NumGoroutine())
}
