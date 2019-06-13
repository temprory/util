package util

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrWorkersTimeout = errors.New("cors timeout")
	ErrWorkersStopped = errors.New("cors stopped")
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

type workerTask struct {
	wg     *sync.WaitGroup
	caller func()
}

type Workers struct {
	sync.WaitGroup
	tag     string
	ch      chan workerTask
	running bool
}

func (c *Workers) runLoop(partition int) {
	Go(func() {
		// logDebug("cors %v child-%v start", c.tag, partition)
		defer func() {
			// /logDebug("cors %v child-%v exit", c.tag, partition)
			c.Done()
		}()
		for task := range c.ch {
			if task.wg == nil {
				Safe(task.caller)
			} else {
				func() {
					defer task.wg.Done()
					defer HandlePanic()
					task.caller()
				}()
			}
			c.Done()
		}
	})
}

func (c *Workers) Go(h func()) error {
	if !c.running {
		return ErrWorkersStopped
	}
	c.Add(1)
	c.ch <- workerTask{nil, h}
	return nil
}

func (c *Workers) GoWithTimeout(h func(), to time.Duration) error {
	if !c.running {
		return ErrWorkersStopped
	}
	c.Add(1)
	select {
	case c.ch <- workerTask{nil, h}:
	case <-time.After(to):
		c.Done()
		return ErrWorkersTimeout
	}
	return nil
}

func (c *Workers) GoWait(h func()) error {
	if !c.running {
		return ErrWorkersStopped
	}
	c.Add(1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c.ch <- workerTask{wg, h}
	wg.Wait()
	return nil
}

func (c *Workers) GoWaitWithTimeout(h func(), to time.Duration) error {
	if !c.running {
		return ErrWorkersStopped
	}
	c.Add(1)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	select {
	case c.ch <- workerTask{wg, h}:
	case <-time.After(to):
		c.Done()
		wg.Done()
		return ErrWorkersTimeout
	}

	wg.Wait()

	return nil
}

func (c *Workers) Stop() {
	c.running = false
	close(c.ch)
	c.Wait()
}

func NewWorkers(tag string, qCap int, corNum int) *Workers {
	c := &Workers{
		tag:     tag,
		ch:      make(chan workerTask, qCap),
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

func (task *LinkTask) Done(data interface{}) {
	task.ch <- data
}

func (task *LinkTask) WaitPre() interface{} {
	if task.pre != nil {
		return <-task.pre
	}
	return nil
}

type WorkersLink struct {
	sync.Mutex
	sync.WaitGroup
	tag     string
	ch      chan *LinkTask
	pre     *LinkTask
	running bool
}

func (c *WorkersLink) runLoop(partition int) {
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

func (c *WorkersLink) Go(h func(task *LinkTask)) error {
	if !c.running {
		return ErrWorkersStopped
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

// func (c *WorkersLink) GoWithTimeout(h func(task *LinkTask), to time.Duration) error {
// 	if !c.running {
// 		return ErrWorkersStopped
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
// 		return ErrWorkersTimeout
// 	}
// 	return nil
// }

func (c *WorkersLink) Stop() {
	c.running = false
	close(c.ch)
	c.Wait()
}

func (c *WorkersLink) StopAsync() {
	Go(func() {
		c.running = false
		close(c.ch)
		c.Wait()
	})
}

func NewWorkersLink(tag string, qCap int, corNum int) *WorkersLink {
	c := &WorkersLink{
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
