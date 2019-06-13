package main

import (
	"fmt"
	"github.com/temprory/util"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("cornum 1: ", runtime.NumGoroutine())

	wg := sync.WaitGroup{}
	cl := util.NewWorkersLink("test", 100, 20)
	for i := 1; i <= 50; i++ {
		idx := i
		wg.Add(1)
		cl.Go(func(task *util.LinkTask) {
			defer wg.Done()

			task.WaitPre()

			fmt.Println("---", idx)
			if idx%20 == 0 {
				fmt.Println("+++ sleep", idx)
				time.Sleep(time.Second)
			}

			task.Done(nil)
		})

	}
	wg.Wait()
	cl.Stop()
	fmt.Println("cornum 2: ", runtime.NumGoroutine())
}
