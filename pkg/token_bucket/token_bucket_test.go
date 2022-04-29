package token_bucket

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTokenBucket_PopToken(t1 *testing.T) {
	tokenBucket := NewTokenBucket(10000, 1000)

	count := 0
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		for time.Now().Sub(start) < time.Millisecond*5400 {
			ok := tokenBucket.PopToken()
			if ok {
				count++
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	wg.Wait()

	fmt.Println(count)
}
