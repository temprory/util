package main

import (
	"fmt"
	"github.com/temprory/util"
	"time"
)

func main() {

	m := map[int]interface{}{}

	for i := 0; i < 4; i++ {
		idx := i
		util.Go(func() {
			for {
				m[idx] = idx
				fmt.Println("-----:", idx)
				time.Sleep(time.Second / 1000)
			}
		})
	}

	time.Sleep(time.Second * 1000000)
}
