package main

import (
	"fmt"
	"github.com/temprory/util"
	"math/rand"
	"sync"
	"time"
)

func main() {
	cors := util.NewCors("test", 10, 5)
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		idx := i
		wg.Add(1)
		cors.Go(func() {
			defer wg.Done()
			time.Sleep(time.Second / time.Duration(10*(1+rand.Intn(10))))
			fmt.Println(idx)
		})
	}
	wg.Wait()
	cors.Stop()
	fmt.Println("over")
}
