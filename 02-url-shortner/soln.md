# Go Standard Library Reference: HTTP Servers & URL Handling

Consolidated reference of the Go standard library packages and methods used in the URL shortener (gophercises) project, along with closely related alternatives.

---

## 1. `net/http` (HTTP Server, Routing, Requests, Responses)

| Method / Symbol | Usage Type | Signature | Purpose |
| :--- | :--- | :--- | :--- |
| **`http.NewServeMux()`** | **Used** | `func() *ServeMux` | Creates a new, empty router. Prefer this over the global default. |
| **`(*ServeMux).HandleFunc(pattern, handler)`** | **Used** | `func(pattern string, handler func(ResponseWriter, *Request))` | Registers a plain function as the handler for a URL pattern, on that specific mux. |
| **`http.ListenAndServe(addr, handler)`** | **Used** | `func(addr string, handler Handler) error` | Starts an HTTP server on `addr`. Blocks until it stops or errors. `nil` handler falls back to the global default mux. |
| **`http.Redirect(w, r, url, code)`** | **Used** | `func(w ResponseWriter, r *Request, url string, code int)` | Sends the client a redirect response — sets `Location` header and status code. |
| **`http.HandlerFunc`** (type) | **Used** | `type HandlerFunc func(ResponseWriter, *Request)` | Adapter type — wraps a plain `func(w, r)` so it satisfies the `Handler` interface. |
| **`http.Handler`** (interface) | **Used** | `ServeHTTP(w ResponseWriter, r *Request)` | Core interface. Anything with this method can serve requests. Used as the type for `fallback` params, mux, server. |
| **`(ResponseWriter).WriteHeader(code)`** | **Used** | `func(statusCode int)` | Locks in and sends the HTTP status code. Must be called before any `Write`. |
| `http.HandleFunc(pattern, handler)` | *Related (avoid)* | `func(pattern string, handler func(ResponseWriter, *Request))` | Package-level version — registers on the global `DefaultServeMux`. Risky: any imported package can silently add routes to it. |
| `http.Get(url)` | *Related* | `func(url string) (*Response, error)` | Client-side GET request using `http.DefaultClient` (no timeout by default). |
| `http.Server` (struct) | *Related* | `struct{ Addr string; Handler Handler; ReadTimeout, WriteTimeout time.Duration }` | Explicit server config — lets you set timeouts. `ListenAndServe` builds one of these internally with no timeouts set. |
| `http.NotFoundHandler()` | *Related* | `func() Handler` | Returns a handler that always replies 404 — useful as a base-case fallback. |

### Relevant status code constants

| Constant | Value | Meaning |
| :--- | :--- | :--- |
| **`http.StatusFound`** | 302 | Temporary redirect — used for the shortener redirect. |
| **`http.StatusNotFound`** | 404 | Used in the custom fallback handler. |
| `http.StatusMovedPermanently` | 301 | Permanent redirect — not used, but common alternative. |

### `*http.Request` fields used

| Field | Type | Purpose |
| :--- | :--- | :--- |
| **`r.URL.Path`** | `string` | The path component of the request, e.g. `/go`. |
| **`r.Header`** | `http.Header` (`map[string][]string`) | Request headers — iterated with `range` in the `headers` handler. |
| `r.Method` | `string` | HTTP method, e.g. `"GET"`. |
| `r.URL.Query()` | `url.Values` (`map[string][]string`) | Query string params — not yet used, relevant for reading things like `?dest=...`. |

---

## 2. `net/url` (used indirectly via `r.URL`)

| Method / Symbol | Usage Type | Signature | Purpose |
| :--- | :--- | :--- | :--- |
| **`url.URL`** (type, accessed via `r.URL`) | **Used** | `struct{ Path string; RawQuery string; ... }` | Backing type of `r.URL`. Not imported directly — accessed through `*http.Request`. |
| `(*url.URL).Query()` | *Related* | `func() Values` | Parses `RawQuery` into `url.Values`. Not used yet in this project. |

---

## 3. `fmt` (Formatting & Console/Response I/O)

| Method / Symbol | Usage Type | Signature | Purpose |
| :--- | :--- | :--- | :--- |
| **`fmt.Fprintf(w, format, a...)`** | **Used** | `func(w io.Writer, format string, a ...any) (int, error)` | Writes formatted output to any `io.Writer` — works with `http.ResponseWriter` since it satisfies that interface. Used to write the HTTP response body. |
| **`fmt.Printf(format, a...)`** | **Used** | `func(format string, a ...any) (int, error)` | Writes formatted output to stdout — used here for server-side debug logging (not sent to the client). |
| **`fmt.Println(a...)`** | **Used** | `func(a ...any) (int, error)` | Writes to stdout with default formatting and a trailing newline. |
| `fmt.Sprintf(format, a...)` | *Related* | `func(format string, a ...any) string` | Same formatting as `Printf`, but returns a string instead of printing. |

---

## 4. `gopkg.in/yaml.v3` (Third-Party: YAML Parsing)

Install: `go get gopkg.in/yaml.v3`

| Method / Symbol | Usage Type | Signature | Purpose |
| :--- | :--- | :--- | :--- |
| **`yaml.Unmarshal(in, out)`** | **Used** | `func(in []byte, out any) error` | Parses YAML bytes into the Go value pointed to by `out` (usually a pointer to a slice or struct). |

### Struct tags

```go
type pathURL struct {
    Path string `yaml:"path"`
    URL  string `yaml:"url"`
}
```
`` `yaml:"..."` `` tags map YAML keys to struct fields explicitly. Without a tag, `Unmarshal` falls back to matching the lowercased field name.

---

## 5. `log` (used for error handling in `main`)

| Method / Symbol | Usage Type | Signature | Purpose |
| :--- | :--- | :--- | :--- |
| **`log.Fatal(v...)`** | **Used** | `func(v ...any)` | Prints the message and exits immediately (`os.Exit(1)`). Used after `YAMLHandler` returns an error, since parsing can fail. |
| `log.Fatalf(format, v...)` | *Related* | `func(format string, v ...any)` | Same as `Fatal`, but with `Printf`-style formatting. |

---

## Quick signature cheat-sheet

| Package | Symbol | Signature |
| :--- | :--- | :--- |
| `net/http` | `NewServeMux` | `func() *ServeMux` |
| `net/http` | `(*ServeMux).HandleFunc` | `func(pattern string, handler func(ResponseWriter, *Request))` |
| `net/http` | `ListenAndServe` | `func(addr string, handler Handler) error` |
| `net/http` | `Redirect` | `func(w ResponseWriter, r *Request, url string, code int)` |
| `net/http` | `HandlerFunc` (type) | `type HandlerFunc func(ResponseWriter, *Request)` |
| `net/http` | `Handler` (interface) | `ServeHTTP(w ResponseWriter, r *Request)` |
| `net/http` | `(ResponseWriter).WriteHeader` | `func(statusCode int)` |
| `fmt` | `Fprintf` | `func(w io.Writer, format string, a ...any) (int, error)` |
| `fmt` | `Printf` | `func(format string, a ...any) (int, error)` |
| `gopkg.in/yaml.v3` | `Unmarshal` | `func(in []byte, out any) error` |
| `log` | `Fatal` | `func(v ...any)` |

---

## Project-defined functions (package `urlshortner`)

| Function | Signature | Purpose |
| :--- | :--- | :--- |
| `MapHandler` | `func(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc` | Returns a handler that redirects based on a lookup map, or delegates to `fallback`. |
| `YAMLHandler` | `func(yml []byte, fallback http.Handler) (http.HandlerFunc, error)` | Parses YAML, builds a map, delegates to `MapHandler`. Returns an error since parsing can fail. |
| `parseYAML` | `func(data []byte) ([]pathURL, error)` | Wraps `yaml.Unmarshal`, returns a typed slice. |
| `buildMap` | `func(pathURLs []pathURL) map[string]string` | Converts `[]pathURL` into the map shape `MapHandler` expects. |