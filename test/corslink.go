package main

import (
	"fmt"
	"github.com/temprory/util"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("cornum 1: ", runtime.NumGoroutine())
	cl := util.NewCorsLink("test", 100, 20)
	time.Sleep(time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < 50; i++ {
		idx := i
		delay := time.Second + time.Duration(rand.Intn(10))*time.Second/10

		cl.Go(func(task *util.LinkTask) {
			//fmt.Println("begin", idx, delay)
			time.Sleep(delay)
			task.Wait()
			fmt.Println("end", idx)
			task.Done(idx)
			if idx == 49 {
				defer wg.Done()
				cl.StopAsync()
				time.Sleep(time.Second)

			}

		})

	}
	wg.Wait()
	fmt.Println("cornum 2: ", runtime.NumGoroutine())
}
