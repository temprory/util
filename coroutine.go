package util

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCorsTimeout = errors.New("cors timeout")
	ErrCorsStopped = errors.New("cors stopped")
)

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
	tag     string
	ch      chan func()
	running bool
}

func (c *Cors) runLoop(partition int) {
	Go(func() {
		// logDebug("cors %v child-%v start", c.tag, partition)
		defer func() {
			// /logDebug("cors %v child-%v exit", c.tag, partition)
			c.Done()
		}()
		for h := range c.ch {
			Safe(h)
			c.Done()
		}
	})
}

func (c *Cors) Go(h func()) error {
	if !c.running {
		return ErrCorsStopped
	}
	c.Add(1)
	c.ch <- h
	return nil
}

func (c *Cors) GoWithTimeout(h func(), to time.Duration) error {
	if !c.running {
		return ErrCorsStopped
	}
	c.Add(1)
	select {
	case c.ch <- h:
	case <-time.After(to):
		c.Done()
		return ErrCorsTimeout
	}
	return nil
}

func (c *Cors) Stop() {
	c.running = false
	close(c.ch)
	c.Wait()
}

func NewCors(tag string, qCap int, corNum int) *Cors {
	c := &Cors{
		tag:     tag,
		ch:      make(chan func(), qCap),
		running: true,
	}
	for i := 0; i < corNum; i++ {
		c.Add(1)
		c.runLoop(i)
	}
	return c
}

type LinkTask struct {
	caller func(*LinkTask)
	pre    chan interface{}
	ch     chan interface{}
}

func (task *LinkTask) Next(data interface{}) {
	task.ch <- data
}

func (task *LinkTask) Pre() interface{} {
	if task.pre != nil {
		return <-task.pre
	}
	return nil
}

type CorsLink struct {
	sync.Mutex
	sync.WaitGroup
	tag     string
	ch      chan *LinkTask
	pre     *LinkTask
	running bool
}

func (c *CorsLink) runLoop(partition int) {
	Go(func() {
		// logDebug("cors %v child-%v start", c.tag, partition)
		defer func() {
			// logDebug("cors %v child-%v exit", c.tag, partition)
			c.Done()
		}()
		for task := range c.ch {
			func() {
				defer HandlePanic()
				task.caller(task)

			}()
			c.Done()
		}
	})
}

func (c *CorsLink) Go(h func(task *LinkTask)) error {
	if !c.running {
		return ErrCorsStopped
	}
	c.Add(1)

	c.Lock()
	var task *LinkTask
	if c.pre != nil {
		task = &LinkTask{
			caller: h,
			pre:    c.pre.ch,
			ch:     make(chan interface{}, 1),
		}
	} else {
		task = &LinkTask{
			caller: h,
			ch:     make(chan interface{}, 1),
		}
	}
	c.pre = task
	c.Unlock()
	c.ch <- task
	return nil
}

// func (c *CorsLink) GoWithTimeout(h func(task *LinkTask), to time.Duration) error {
// 	if !c.running {
// 		return ErrCorsStopped
// 	}
// 	c.Add(1)
// 	c.Lock()
// 	defer c.Unlock()

// 	var task *LinkTask
// 	if c.pre != nil {
// 		task = &LinkTask{
// 			caller: h,
// 			pre:    c.pre.ch,
// 			ch:     make(chan interface{}),
// 		}
// 	} else {
// 		task = &LinkTask{
// 			caller: h,
// 			ch:     make(chan interface{}),
// 		}
// 	}

// 	select {
// 	case c.ch <- task:
// 		c.pre = task
// 	case <-time.After(to):
// 		c.Done()
// 		return ErrCorsTimeout
// 	}
// 	return nil
// }

func (c *CorsLink) Stop() {
	c.running = false
	close(c.ch)
	c.Wait()
}

func (c *CorsLink) StopAsync() {
	Go(func() {
		c.running = false
		close(c.ch)
		c.Wait()
	})
}

func NewCorsLink(tag string, qCap int, corNum int) *CorsLink {
	c := &CorsLink{
		tag:     tag,
		ch:      make(chan *LinkTask, qCap),
		running: true,
	}
	for i := 0; i < corNum; i++ {
		c.Add(1)
		c.runLoop(i)
	}

	return c
}
