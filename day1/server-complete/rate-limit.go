package main

import (
	"sync"
	"time"
)

type rateLimitStruct struct {
	m       *sync.Mutex
	records map[string]rateLimitRecordStruct
	max     int
}

func newRateLimit(max int) *rateLimitStruct {
	rateLimit := &rateLimitStruct{
		m:       &sync.Mutex{},
		records: map[string]rateLimitRecordStruct{},
		max:     max,
	}
	return rateLimit
}

type rateLimitRecordStruct struct {
	count  int
	window int64
}

func (rateLimiter *rateLimitStruct) check(key string) bool {
	rateLimiter.m.Lock()
	defer rateLimiter.m.Unlock()

	window := time.Now().Unix() / 60
	record, ok := rateLimiter.records[key]
	if !ok || record.window != window {
		rateLimiter.records[key] = rateLimitRecordStruct{
			count:  1,
			window: window,
		}
		return true
	}
	if record.count >= rateLimiter.max {
		return false
	}
	record.count++
	rateLimiter.records[key] = record
	return true
}
