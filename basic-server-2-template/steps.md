# Steps

## GET /

HTML を表示する。

```go
if r.Method == "GET" && r.URL.Path == "/" {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("<h1>hello</h1>"))
	return
}
```

## GET /google

`google.com` にリダイレクト。

```go
if r.Method == "GET" && r.URL.Path == "/google" {
	w.Header().Set("Location", "https://google.com")
	w.WriteHeader(303)
	return
}
```

## POST /encode

1. リクエストボディを読み取る。
2. リクエストボディを base64 変換。
3. base64 の文字列をを返す。

```go
import (
	"encoding/base64"
	"io"
)
```

```go
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
```

## POST /add

```
go get github.com/faroedev/go-json
```

```go
import (
	"github.com/faroedev/go-json"
	"io"
)
```

1. リクエストボディを読み取る。
2. JSON オブジェクトして parse する。
3. `a` の値（`int`）を読み取る。
4. `b` の値（`int`）を読み取る。
5. `a + b` を `sum` の値として JSON オブジェクトを作成。
6. JSON を返す。


```go
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
	resultJSONBuilder.AddInt("result", a+b)
	resultJSON := resultJSONBuilder.Done()
	w.Write([]byte(resultJSON))
	return
}
```