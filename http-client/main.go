package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/faroedev/go-json"
)

//go:embed page.html
var pageHTML []byte

func main() {
	port := 3000
	if len(os.Args) > 1 {
		parsed, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Invalid port number '%s'\n", os.Args[1])
		}
		port = parsed
	}
	fmt.Printf("Starting server on port %d...\n", port)
	address := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(address, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(pageHTML)
			return
		}
		if r.Method == "POST" && r.URL.Path == "/" {
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				writeErrorResponseBody(w, "Invalid request body.")
				return
			}
			bodyJSONObject, err := json.ParseObject(string(requestBody))
			if err != nil {
				w.WriteHeader(400)
				writeErrorResponseBody(w, "Invalid request body.")
				return
			}
			method, err := bodyJSONObject.GetString("method")
			if err != nil {
				w.WriteHeader(400)
				writeErrorResponseBody(w, "Invalid or missing 'method' field.")
				return
			}
			uri, err := bodyJSONObject.GetString("uri")
			if err != nil {
				w.WriteHeader(400)
				writeErrorResponseBody(w, "Invalid or missing 'uri' field.")
				return
			}
			body, err := bodyJSONObject.GetString("body")
			if err != nil {
				w.WriteHeader(400)
				writeErrorResponseBody(w, "Invalid or missing 'body' field.")
				return
			}
			headerFieldsJSONArray, err := bodyJSONObject.GetJSONArray("header_fields")
			if err != nil {
				w.WriteHeader(400)
				writeErrorResponseBody(w, "Invalid or missing 'header_fields' field.")
				return
			}
			if headerFieldsJSONArray.Length()%2 != 0 {
				w.WriteHeader(400)
				writeErrorResponseBody(w, "Invalid or missing 'header_fields' field.")
				return
			}
			headerFields := []string{}
			for i := range headerFieldsJSONArray.Length() {
				headerItem, err := headerFieldsJSONArray.GetString(i)
				if err != nil {
					w.WriteHeader(400)
					writeErrorResponseBody(w, "Invalid or missing 'header_fields' field.")
					return
				}
				headerFields = append(headerFields, headerItem)
			}

			httpRequest, err := http.NewRequest(method, uri, strings.NewReader(body))
			if err != nil {
				w.WriteHeader(400)
				writeErrorResponseBody(w, fmt.Sprintf("Failed to create request: %s.", err.Error()))
				return
			}
			for i := 0; i < len(headerFields); i += 2 {
				httpRequest.Header.Add(headerFields[i], headerFields[i+1])
			}
			httpResponse, err := http.DefaultClient.Do(httpRequest)
			if err != nil {
				w.WriteHeader(500)
				writeErrorResponseBody(w, fmt.Sprintf("Failed to send request: %s.", err.Error()))
				return
			}
			responseBody, err := io.ReadAll(httpResponse.Body)
			if err != nil {
				w.WriteHeader(500)
				writeErrorResponseBody(w, fmt.Sprintf("Failed to read response body: %s.", err.Error()))
				return
			}
			resultHeaderFieldsJSONBuilder := json.NewArrayBuilder()
			for fieldName, fieldValues := range httpResponse.Header {
				for _, fieldValue := range fieldValues {
					resultHeaderFieldsJSONBuilder.AddString(fieldName)
					resultHeaderFieldsJSONBuilder.AddString(fieldValue)
				}
			}
			resultHeaderFieldsJSON := resultHeaderFieldsJSONBuilder.Done()

			resultJSONBuilder := json.NewObjectBuilder()
			resultJSONBuilder.AddInt("status", httpResponse.StatusCode)
			resultJSONBuilder.AddJSON("header_fields", resultHeaderFieldsJSON)
			resultJSONBuilder.AddString("body", string(responseBody))
			resultJSON := resultJSONBuilder.Done()

			w.Write([]byte(resultJSON))
			return
		}
		w.WriteHeader(404)
	}))
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
}

func writeErrorResponseBody(w io.Writer, message string) {
	builder := json.NewObjectBuilder()
	builder.AddString("message", message)
	resultJSON := builder.Done()
	w.Write([]byte(resultJSON))
}
