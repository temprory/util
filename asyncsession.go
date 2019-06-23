package util

import (
	"fmt"
	"sync"
	"time"
)

type SessionMgr struct {
	sync.Mutex
	sessions map[interface{}]chan interface{}
}

func (mgr *SessionMgr) Add(sess interface{}) error {

	mgr.Lock()
	if _, ok := mgr.sessions[sess]; ok {
		mgr.Unlock()
		return fmt.Errorf("session %v exist", sess)
	}
	mgr.sessions[sess] = make(chan interface{}, 1)
	mgr.Unlock()

	return nil
}

func (mgr *SessionMgr) Done(sess interface{}, data interface{}) error {
	mgr.Lock()
	if done, ok := mgr.sessions[sess]; ok {
		//delete(mgr.sessions, sess)
		mgr.Unlock()
		done <- data
		return nil
	}
	mgr.Unlock()

	return fmt.Errorf("session %v not exist", sess)
}

func (mgr *SessionMgr) Wait(sess interface{}, timeout time.Duration) (interface{}, error) {
	mgr.Lock()
	done, ok := mgr.sessions[sess]
	if ok {
		mgr.Unlock()

		var data interface{}
		if timeout > 0 {
			select {
			case data = <-done:
			case <-time.After(timeout):
				mgr.Lock()
				delete(mgr.sessions, sess)
				mgr.Unlock()
				return nil, fmt.Errorf("wait session %v timeout", sess)
			}
		} else {
			data = <-done
		}

		mgr.Lock()
		delete(mgr.sessions, sess)
		mgr.Unlock()

		return data, nil
	}

	mgr.Unlock()

	return nil, fmt.Errorf("session %v not exist", sess)
}

func (mgr *SessionMgr) Len() int {
	mgr.Lock()
	l := len(mgr.sessions)
	mgr.Unlock()
	return l
}

func NewSessionMgr() *SessionMgr {
	return &SessionMgr{
		sessions: map[interface{}]chan interface{}{},
	}
}
