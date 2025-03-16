//go:build rate_limiter

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const RATE_LIMIT = 5
const PERIOD = 1 * time.Second

type RateLimiter interface {
	AllowRequest(userID string) bool         // Returns true if request is allowed, false otherwise
	Service(chan os.Signal, *sync.WaitGroup) // Cleans user rates per period
}

type MyRateLimiter struct {
	default_limit     int                  // Allowed number of requests
	tokens_per_period map[string]int       // Tokens per user
	last_refill       map[string]time.Time // Last token quota re-fill timestamp
	period            time.Duration        // Requests limit period
	m                 sync.Mutex
}

func (rl *MyRateLimiter) AllowRequest(userID string) bool {
	rl.m.Lock()
	defer rl.m.Unlock()

	// initiate if user does not exists, yet
	if _, exists := rl.tokens_per_period[userID]; !exists {
		rl.tokens_per_period[userID] = rl.default_limit
		rl.last_refill[userID] = time.Now()
	}

	fmt.Printf("%s: count: %d\n", userID, rl.tokens_per_period[userID])

	if rl.tokens_per_period[userID] > 0 {
		rl.tokens_per_period[userID]--
		return true
	}
	return false
}

func (rl *MyRateLimiter) Service(quit chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(rl.period)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Reseting request counts ..")
			rl.m.Lock()
			now := time.Now()
			for userID, last_refill := range rl.last_refill {
				if last_refill.Add(2 * rl.period).Before(now) { // Remove inactive users
					delete(rl.tokens_per_period, userID)
					delete(rl.last_refill, userID)
				} else {
					rl.last_refill[userID] = time.Now()
					rl.tokens_per_period[userID] = rl.default_limit
				}
			}
			rl.m.Unlock()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func NewRateLimiter(limit int, period time.Duration) RateLimiter {
	return &MyRateLimiter{
		default_limit:     limit,
		tokens_per_period: make(map[string]int),
		last_refill:       make(map[string]time.Time),
		period:            period,
		m:                 sync.Mutex{},
	}
}

func GenerateTraffic(rl RateLimiter, userID string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 100; i++ {
		if rl.AllowRequest(userID) {
			fmt.Printf("%s: request allowed\n", userID)
		} else {
			fmt.Printf("%s: rate limit exceeded\n", userID)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup
	quit := make(chan os.Signal, 1)
	rl := NewRateLimiter(RATE_LIMIT, PERIOD)

	wg.Add(1)
	go rl.Service(quit, &wg)

	// make some generated requests
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go GenerateTraffic(rl, fmt.Sprintf("user#%d", i), &wg)
	}

	// shut down listener
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT)

	close(quit)
	wg.Wait()
}
