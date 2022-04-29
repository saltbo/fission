package token_bucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	cap     int64
	rate    int64
	leftCap int64
	ticker  *time.Ticker
	mutex   sync.Mutex
}

func NewTokenBucket(cap, rate int64) *TokenBucket {
	tokenBucket := &TokenBucket{
		cap:     cap,
		rate:    rate,
		leftCap: 0,
	}

	tokenBucket.ticker = time.NewTicker(time.Second)
	incrPreTicker := rate
	if rate >= 10 {
		tokenBucket.ticker = time.NewTicker(time.Duration(100) * time.Millisecond)
		incrPreTicker = rate / 10
	}

	go func() {
		for {
			if _, ok := <-tokenBucket.ticker.C; !ok {
				return
			}

			tokenBucket.mutex.Lock()
			if tokenBucket.leftCap+incrPreTicker <= tokenBucket.cap {
				tokenBucket.leftCap = tokenBucket.leftCap + incrPreTicker
			} else {
				tokenBucket.leftCap = tokenBucket.cap
			}
			tokenBucket.mutex.Unlock()
		}
	}()

	return tokenBucket
}

func (t *TokenBucket) PopToken() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.leftCap > 0 {
		t.leftCap--
		return true
	}
	return false
}

func (t *TokenBucket) Stop() {
	t.ticker.Stop()
}
