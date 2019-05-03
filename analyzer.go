package util

import (
	"fmt"
	"github.com/temprory/log"
	"runtime"
	"strings"
	"time"
)

var (
	analyzerDebug = false
)

func SetAnalyzerDebug(d bool) {
	analyzerDebug = d
}

type Analyzer struct {
	Tag        string
	Parent     *Analyzer `json:"-"`
	Children   []*Analyzer
	Limit      time.Duration
	TBegin     time.Time
	TEnd       time.Time
	StackBegin string
	StackEnd   string
	Data       interface{}
	expired    bool
}

func getStackInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = -1
	} else {
		if log.BuildDir != "" {
			if strings.HasPrefix(file, log.BuildDir) {
				file = file[len(log.BuildDir):]
			}
		}
	}
	return fmt.Sprintf("[%s %d]", file, line)
}

func (a *Analyzer) Begin() {
	a.StackBegin = getStackInfo()
}

func (a *Analyzer) Done(v ...interface{}) {
	a.TEnd = time.Now()
	a.StackEnd = getStackInfo()
	if len(v) > 0 {
		a.Data = v[0]
	}
	a.expired = a.TEnd.Sub(a.TBegin) > a.Limit
	if a.expired {
		tmp := a.Parent
		for tmp != nil {
			//fmt.Println(a.Tag)
			tmp.expired = true
			tmp = tmp.Parent
		}
	}
}

func (a *Analyzer) Fork(tag string, limit time.Duration) *Analyzer {
	analyzer := &Analyzer{
		Tag:    tag,
		Parent: a,
		Limit:  limit,
		TBegin: time.Now(),
	}
	a.Children = append(a.Children, analyzer)
	return analyzer
}

// func (a *Analyzer) Info(args ...interface{}) (string, bool) {
// 	indent := ""
// 	if len(args) > 0 {
// 		indent = args[0].(string)
// 	}

// 	used := a.TEnd.Sub(a.TBegin)
// 	expired := used > a.Limit
// 	infoStr := fmt.Sprintf("%v[%v] cost: %vus, exp: %v\n", indent, a.Tag, used.Nanoseconds()/1000, expired)
// 	indent += "--"
// 	for _, v := range a.Children {
// 		str, exp := v.Info(indent)
// 		if exp {
// 			expired = exp
// 		}
// 		infoStr += str
// 	}
// 	return infoStr, expired
// }

func (a *Analyzer) Expired() bool {
	// used := a.TEnd.Sub(a.TBegin)
	// expired := used > a.Limit
	// for _, v := range a.Children {
	// 	if v.Expired() {
	// 		return true
	// 	}
	// }
	return a.expired
}

func (a *Analyzer) Info() string {
	if analyzerDebug {
		data, _ := json.MarshalIndent(a, "", "    ")
		return string(data)
	}
	str, _ := json.MarshalToString(a)
	return str
}

func NewAnalyzer(tag string, limit time.Duration) *Analyzer {
	return &Analyzer{
		Tag:    tag,
		Limit:  limit,
		TBegin: time.Now(),
	}
}
