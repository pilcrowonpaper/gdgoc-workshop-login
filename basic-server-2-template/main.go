package main

import (
	"net/http"
)

func main() {
	err := http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/" {
			// GET /
		}

		if r.Method == "GET" && r.URL.Path == "/google" {
			// GET /google
		}

		if r.Method == "POST" && r.URL.Path == "/encode" {
			// POST /encode
		}

		if r.Method == "POST" && r.URL.Path == "/add" {
			// POST /add
		}

		w.WriteHeader(404)
		w.Write([]byte("Not found."))
	}))

	if err != nil {
		panic(err)
	}
}
