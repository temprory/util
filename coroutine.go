package util

import (
	"errors"
	"sync"
	"time"
)

var ErrTimeOut = errors.New("timeout")

func Go(cb func()) {
	go func() {
		defer HandlePanic()
		cb()
	}()
}

// func Go(cb interface{}, args ...interface{}) {
// 	f := reflect.ValueOf(cb)
// 	go func() {
// 		defer HandlePanic(true)
// 		n := len(args)
// 		if n > 0 {
// 			refargs := make([]reflect.Value, n)
// 			for i := 0; i < n; i++ {
// 				refargs[i] = reflect.ValueOf(args[i])
// 			}
// 			f.Call(refargs)

// 		} else {
// 			f.Call(nil)
// 		}
// 	}()
// }

type Cors struct {
	sync.WaitGroup
	tag string
	ch  chan func()
}

func (c *Cors) runLoop(partition int) {
	Go(func() {
		logDebug("cors %v child-%v start", c.tag, partition)
		defer func() {
			logDebug("cors %v child-%v exit", c.tag, partition)
			c.Done()
		}()
		for h := range c.ch {
			Safe(h)
			c.Done()
		}
	})
}

func (c *Cors) Go(h func()) {
	c.Add(1)
	c.ch <- h
}

func (c *Cors) GoWithTimeout(h func(), to time.Duration) error {
	c.Add(1)
	select {
	case c.ch <- h:
	case <-time.After(to):
		c.Done()
		return ErrTimeOut
	}
	return nil
}

func (c *Cors) Stop() {
	close(c.ch)
	c.Wait()
}

func NewCors(tag string, qCap int, corNum int) *Cors {
	c := &Cors{
		tag: tag,
		ch:  make(chan func(), qCap),
	}
	for i := 0; i < corNum; i++ {
		c.Add(1)
		c.runLoop(i)
	}
	return c
}
