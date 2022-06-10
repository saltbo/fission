package token_bucket

import (
	"sync"
	"time"
)

var mCache sync.Map
var gMutex sync.Mutex

type TimeoutMutex struct {
	time time.Time
	ch   chan struct{}
}

func LockOrWait(name string, callback func()) bool {
	var tMu *TimeoutMutex
	gMutex.Lock()
	val, isLoad := mCache.Load(name)
	if !isLoad {
		tMu = &TimeoutMutex{
			time: time.Now(),
			ch:   make(chan struct{}),
		}
		mCache.Store(name, tMu)
	}
	gMutex.Unlock()

	if isLoad {
		oldMu, ok := val.(*TimeoutMutex)
		if ok {
			<-oldMu.ch
		}
		return false
	} else {
		var runOver bool

		go func() {
			time.Sleep(time.Second * 2)
			if !runOver {
				mCache.Delete(name)
				runOver = true
				close(tMu.ch)
			}
		}()

		callback()
		if !runOver {
			mCache.Delete(name)
			runOver = true
			close(tMu.ch)
		}
		return true
	}
}
