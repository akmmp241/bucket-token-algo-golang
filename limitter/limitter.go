package limitter

import (
	"log"
	"sync"
	"time"
)

type TokenBucket struct {
	Capacity     int
	CurrentToken int
	Rate         int
	Last         time.Time
	MU           sync.Mutex
	RefillAmount int
}

func (tb *TokenBucket) TakeTokens(tokens int) bool {
	tb.MU.Lock()
	defer tb.MU.Unlock()
	tb.refillTokens()

	if tb.CurrentToken >= tokens {
		tb.CurrentToken -= tokens
		log.Println("Current Token: ", tb.CurrentToken)
		return true
	}
	return false
}

func (tb *TokenBucket) refillTokens() {
	now := time.Now()
	diff := now.Sub(tb.Last)
	tb.Last = now

	addedTokens := int(diff.Seconds()) * tb.RefillAmount

	tb.CurrentToken += addedTokens
	if tb.CurrentToken > tb.Capacity {
		tb.CurrentToken = tb.Capacity
	}
}

func NewTokenBucket(capacity int, rate int) *TokenBucket {
	return &TokenBucket{
		Capacity:     capacity,
		CurrentToken: capacity,
		Rate:         rate,
		Last:         time.Now(),
		RefillAmount: rate,
	}
}
