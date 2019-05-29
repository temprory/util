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
	for i := 0; i < 20; i++ {
		idx := i
		wg.Add(1)
		cors.Go(func() {
			defer wg.Done()
			time.Sleep(time.Second / time.Duration(10*(1+rand.Intn(10))))
			fmt.Println(idx)
		})
	}
	wg.Wait()

	fmt.Println("-----")
	for i := 0; i < 20; i++ {
		idx := i
		go func() {
			wg.Add(1)
			defer wg.Done()
			err := cors.GoWait(func() {
				time.Sleep(time.Second / 2)
			})
			fmt.Println(idx, err)
		}()
	}
	wg.Wait()

	fmt.Println("-----")
	for i := 0; i < 20; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			err := cors.GoWaitWithTimeout(func() {
				time.Sleep(time.Second / 5 * 6)
			}, time.Second)
			fmt.Println(idx, err)
		}()
	}
	wg.Wait()
	cors.Stop()
	fmt.Println("over")
}
