package token_bucket

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLockOrWait(t *testing.T) {
	var count int64 = 0

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		go func(ii int) {
			LockOrWait("fun1", func() {
				t.Logf("start create func1 %d", ii)
				time.Sleep(time.Second * 2)
				atomic.AddInt64(&count, 1)
				t.Logf("create func1 success %d", ii)
			})
		}(i)
		time.Sleep(time.Second)
	}

	fmt.Printf("count %d \n", count)
	wg.Wait()

}
