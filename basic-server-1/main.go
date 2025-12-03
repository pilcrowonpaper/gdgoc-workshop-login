package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		fmt.Println(r.URL.RequestURI())
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(body))
		w.WriteHeader(200)
	}))
	if err != nil {
		panic(err)
	}
}
