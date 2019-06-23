package main

import (
	"fmt"
	"github.com/temprory/util"
	"time"
)

func main() {
	mgr := util.NewSessionMgr()

	mgr.Add(1)

	util.Go(func() {
		time.Sleep(time.Second / 1000)
		mgr.Done(1, 2)
	})

	data, err := mgr.Wait(1, time.Second*3)
	fmt.Println("---- ", data, err, mgr.Len())

	err = mgr.Add(1)
	fmt.Println("--- 111:", err)

	err = mgr.Add(1)
	fmt.Println("--- 222:", err, mgr.Len())

	err = mgr.Add(2)
	fmt.Println("--- 333:", err, mgr.Len())

	data, err = mgr.Wait(1, time.Second/10)
	fmt.Println("--- 444:", data, err, mgr.Len())

	data, err = mgr.Wait(2, time.Second/10)
	fmt.Println("--- 555:", data, err, mgr.Len())

}
