package main

import (
	"io"
	"net/http"

	"github.com/faroedev/go-json"
	"golang.org/x/sync/semaphore"
	"zombiezen.com/go/sqlite"

	_ "embed"
)

//go:embed page.html
var pageHTML []byte

type serverStruct struct {
	conn                          *sqlite.Conn
	passwordHashingSemaphore      *semaphore.Weighted
	passwordVerificationRateLimit *rateLimitStruct
}

func (server *serverStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(pageHTML)
		return
	}
	if r.Method == "POST" && r.URL.Path == "/action" {
		body := http.MaxBytesReader(nil, r.Body, 1024*1024*16)
		bodyBytes, err := io.ReadAll(body)
		if _, ok := err.(*http.MaxBytesError); ok {
			w.WriteHeader(413)
			return
		}
		if err != nil {
			w.WriteHeader(400)
			return
		}

		bodyJSONObject, err := json.ParseObject(string(bodyBytes))
		if err != nil {
			w.WriteHeader(400)
			return
		}
		action, err := bodyJSONObject.GetString("action")
		if err != nil {
			w.WriteHeader(400)
			return
		}
		argumentsJSONObject, err := bodyJSONObject.GetJSONObject("arguments")
		if err != nil {
			w.WriteHeader(400)
			return
		}

		result, err := server.invokeAction(action, argumentsJSONObject)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(result))
		return
	}

	w.WriteHeader(404)
}
