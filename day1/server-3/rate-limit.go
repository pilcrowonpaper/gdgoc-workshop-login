package main

import (
	"sync"
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

	// TODO
}
