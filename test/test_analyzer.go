package main

import (
	"fmt"
	"github.com/temprory/util"
	"time"
)

func A() {
	a := util.NewAnalyzer("A", time.Second)
	a.Begin()
	defer func() {
		a.Done()
		if a.Expired {
			fmt.Println("expired:", a.Info())
		} else if a.ChildExpired {
			fmt.Println("child expired:", a.Info())
		}
	}()

	B(a.Fork("B", time.Second/20))

	C(a.Fork("C", time.Second/20))
}

func B(a *util.Analyzer) {
	a.Begin()

	time.Sleep(time.Second / 20)

	D(a.Fork("D", time.Second/20))

	a.Done("binfo")
}

func C(a *util.Analyzer) {
	a.Begin()

	D(a.Fork("D", time.Second/20))

	a.Done()
}

func D(a *util.Analyzer) {
	a.Begin()

	time.Sleep(time.Second / 10)

	a.Done()
}

// go build -ldflags "-X github.com/temprory/log.BuildDir=C:/Users/User/Desktop/" test_analyzer.go
func main() {
	util.SetAnalyzerDebug(true)
	A()

	var a *util.Analyzer
	a.Begin()
	a.Done()
	fmt.Println("empty info:", a.Info())
}
