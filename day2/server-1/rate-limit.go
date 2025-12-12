package main

import (
	"sync"
	"time"
)

type rateLimitStruct struct {
	m          *sync.Mutex
	buckets    map[string]tokenBucketStruct
	max        int
	refillRate time.Duration
}

func newRateLimit(max int, refillRate time.Duration) *rateLimitStruct {
	rateLimit := &rateLimitStruct{
		m:          &sync.Mutex{},
		buckets:    map[string]tokenBucketStruct{},
		max:        max,
		refillRate: refillRate,
	}
	return rateLimit
}

type tokenBucketStruct struct {
	count          int
	lastRefilledAt time.Time
}

func (rateLimiter *rateLimitStruct) check(key string) bool {
	rateLimiter.m.Lock()
	defer rateLimiter.m.Unlock()

	now := time.Now()
	bucket, ok := rateLimiter.buckets[key]
	if !ok {
		rateLimiter.buckets[key] = tokenBucketStruct{
			count:          rateLimiter.max - 1,
			lastRefilledAt: now,
		}
		return true
	}
	refill := int(now.Sub(bucket.lastRefilledAt) / rateLimiter.refillRate)
	bucket.count += refill
	bucket.lastRefilledAt = bucket.lastRefilledAt.Add(rateLimiter.refillRate * time.Duration(refill))
	if bucket.count < 1 {
		rateLimiter.buckets[key] = bucket
		return false
	}
	bucket.count--
	rateLimiter.buckets[key] = bucket
	return true
}
