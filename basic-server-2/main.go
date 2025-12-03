package main

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/faroedev/go-json"
)

func main() {
	err := http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(200)
			w.Write([]byte("<h1>hello</h1>"))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/google" {
			w.Header().Set("Location", "https://google.com")
			w.WriteHeader(303)
			return
		}

		if r.Method == "POST" && r.URL.Path == "/encode" {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				return
			}
			encoded := base64.StdEncoding.EncodeToString(b)
			w.Write([]byte(encoded))
			return
		}

		if r.Method == "POST" && r.URL.Path == "/add" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			bodyJSONObject, err := json.ParseObject(string(bodyBytes))
			if err != nil {
				w.WriteHeader(400)
				return
			}

			a, err := bodyJSONObject.GetInt("a")
			if err != nil {
				w.WriteHeader(400)
				return
			}
			b, err := bodyJSONObject.GetInt("b")
			if err != nil {
				w.WriteHeader(400)
				return
			}

			resultJSONBuilder := json.NewObjectBuilder()
			resultJSONBuilder.AddInt("sum", a+b)
			resultJSON := resultJSONBuilder.Done()
			w.Write([]byte(resultJSON))
			return
		}

		w.WriteHeader(404)
		w.Write([]byte("Not found."))
	}))

	if err != nil {
		panic(err)
	}
}
