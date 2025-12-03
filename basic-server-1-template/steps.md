# Steps

1. メソッドを出力
2. URI を出力
3. 各ヘッダーを出力
4. ボディを読み取る
5. ボディを出力
6. ステータス 200 を返す

```go
import (
	"net/http"
)
```

```go
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
```
