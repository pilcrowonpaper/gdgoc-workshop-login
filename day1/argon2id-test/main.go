package main

import (
	"crypto/rand"
	"sync"

	"golang.org/x/crypto/argon2"
)

func main() {
	wg := sync.WaitGroup{}
	for range 1000 {
		go func() {
			salt := make([]byte, 32)
			rand.Read(salt)
			_ = argon2.IDKey([]byte("hello world!"), salt, 3, 64*1024, 1, 32)
			wg.Done()
		}()
		wg.Add(1)
	}
	wg.Wait()
}
