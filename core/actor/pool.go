package actor

import (
	"sync"
	"time"
)

/*
   @Author: orbit-w
   @File: pool
   @2024 2月 周日 22:29
*/

var (
	reqPool = sync.Pool{New: func() any {
		return &Request{
			ch: make(chan *Response, 1),
		}
	}}

	timerPool = sync.Pool{}
)

func GetRequest() *Request {
	v := reqPool.Get()
	return v.(*Request)
}

func AcquireTimer(d time.Duration) *time.Timer {
	v := timerPool.Get()
	if v == nil {
		return time.NewTimer(d)
	}
	t := v.(*time.Timer)
	if t.Reset(d) {
		t = time.NewTimer(d)
	}
	return t
}

func ReleaseTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	timerPool.Put(t)
}
