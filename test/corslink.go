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
			//fmt.Println("begin", idx, delay)
			time.Sleep(delay)
			preData := task.Pre()
			fmt.Println("preData:", preData)
			fmt.Println("end", idx)
			currData := idx
			task.Next(currData)

			allData = append(allData, currData)
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
