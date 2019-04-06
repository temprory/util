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

	allData := []interface{}{}
	for i := 0; i < 50; i++ {
		idx := i
		delay := time.Second + time.Duration(rand.Intn(10))*time.Second/10

		cl.Go(func(task *util.LinkTask) {
			//slow fuck things
			time.Sleep(delay)
			currData := idx

			//wait for pre task done
			preData := task.WaitPre()

			//handle result after slow fuck
			fmt.Println("preData:", preData)
			fmt.Println("end", idx)
			allData = append(allData, currData)

			//release this task
			//task.Done(currData)
			task.Done(nil)

			//退出任务，可选
			if idx == 49 {
				defer wg.Done()
				cl.StopAsync()
				time.Sleep(time.Second)
			}
		})

	}
	wg.Wait()
	for i, v := range allData {
		fmt.Println("--- ", i, v)
	}
	fmt.Println("cornum 2: ", runtime.NumGoroutine())
}
