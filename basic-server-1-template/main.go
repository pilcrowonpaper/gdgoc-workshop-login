package main

import (
	"net/http"
)

func main() {
	err := http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ...
	}))
	if err != nil {
		panic(err)
	}
}
