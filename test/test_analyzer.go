package main

import (
	"fmt"
	"github.com/temprory/util"
	"time"
)

func A() {
	analyzer := util.NewAnalyzer("A", time.Second)
	analyzer.Begin()
	B(analyzer.Fork("B", time.Second/2))
	C(analyzer.Fork("C", time.Second/2))

	aInfo := map[string]interface{}{
		"name":   "test",
		"passwd": "123qwe",
	}
	analyzer.Done(aInfo)

	if analyzer.Expired() {
		fmt.Println(analyzer.Info())
	}
}

func B(analyzer *util.Analyzer) {
	analyzer.Begin()

	time.Sleep(time.Second / 2)

	D(analyzer.Fork("D", time.Second/2))

	analyzer.Done("binfo")
}

func C(analyzer *util.Analyzer) {
	analyzer.Begin()

	D(analyzer.Fork("D", time.Second/2))

	analyzer.Done()
}

func D(analyzer *util.Analyzer) {
	analyzer.Begin()

	time.Sleep(time.Second / 10 * 6)

	analyzer.Done()
}

// go build -ldflags "-X github.com/temprory/log.BuildDir=C:/Users/User/Desktop/" test_analyzer.go
func main() {
	util.SetAnalyzerDebug(true)
	A()
}
