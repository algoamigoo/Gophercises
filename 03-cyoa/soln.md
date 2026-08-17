# Go & Choose Your Own Adventure (CYOA) Learnings

A reference summary of key Go concepts, standard library patterns, and web development fundamentals learned throughout this project.

---

## 1. Struct Tags & Dynamic JSON Mapping

### Dynamic Map Keys

When the root of a JSON file uses arbitrary string keys (e.g., `"intro"`, `"new-york"`, `"denver"`), map it using a Go map type rather than a struct:

```go
type Book map[string]Chapter

```

### Struct Casing & Field Exporting

In Go, fields must be **capitalized** (exported) to be accessible by external packages like `encoding/json` and `html/template`. Struct tags map Go's `PascalCase` fields to JSON's `snake_case` or `camelCase` keys:

```go
type Chapter struct {
    Title   string   `json:"title"`
    Story   []string `json:"story"`
    Options []Option `json:"options"`
}

type Option struct {
    Text string `json:"text"`
    Arc  string `json:"arc"`
}

```

### `omitempty` Behavior

* **Decoding (`json.Unmarshal` / `Decode`):** Struct tags like `,omitempty` have **no effect** when reading JSON into Go.
* **Encoding (`json.Marshal` / `Encode`):** `,omitempty` strips out fields containing zero values (`""`, `nil`, `0`, `[]`) when converting Go structs back into JSON bytes.

---

## 2. Reading Data: `json.Unmarshal` vs. `json.NewDecoder`

| Feature | `json.Unmarshal` | `json.NewDecoder` |
| --- | --- | --- |
| **Input Source** | In-memory `[]byte` slice | Streaming `io.Reader` (`*os.File`, `http.Request.Body`) |
| **Memory Allocation** | Loads the entire file payload into RAM first | Reads data chunk-by-chunk directly from the stream |
| **Best Practice** | Small payloads already loaded in memory | Files on disk, network streams, and HTTP request bodies |

### Decoupling with `io.Reader`

Accepting an `io.Reader` interface instead of a specific file path makes parser functions flexible across files, buffers, and network streams:

```go
func JsonStory(r io.Reader) (Book, error) {
    var book Book
    err := json.NewDecoder(r).Decode(&book)
    if err != nil {
        return nil, err
    }
    return book, nil
}

```

---

## 3. Web Development with `net/http`

### The `http.Handler` Interface

Any type implementing `ServeHTTP(w http.ResponseWriter, r *http.Request)` satisfies the `http.Handler` interface.

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

```

### Dependency Injection via Closures

Since standard HTTP handler signatures cannot accept custom parameters like `Book`, wrap the handler in a constructor function (`NewHandler(b Book)`). The inner function captures the `book` instance via a **closure**:

```go
func NewHandler(b Book) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // The closure has access to 'b' (the Book map) for chapter lookups
        path := r.URL.Path
        // ...
    })
}

```

### URL Path Normalization

Requests arrive with a leading slash (`/denver`). Normalizing paths before map lookups prevents missing keys and handles fallback routes:

1. **Root path fallback:** Convert `"/"` or `""` to `"/intro"`.
2. **Trim leading slash:** Slice `path[1:]` (e.g., `"/denver"` $\rightarrow$ `"denver"`).
3. **Map lookup:** Retrieve `b[path]`. If found, render HTML; if not, return `http.StatusNotFound`.

---

## 4. Templating with `html/template`

Go's `html/template` package parses HTML strings or files and populates them with dynamic Go data types:

```go
var tmpl = template.Must(template.New("").Parse(defaultChapterTemplate))

// Execute writes rendered HTML directly to the http.ResponseWriter
err := tmpl.Execute(w, chapter)

```

### Template Directives Reference

* `{{.Title}}` – Accesses the `Title` field on the current `Chapter` struct.
* `{{range .Story}}<p>{{.}}</p>{{end}}` – Loops over the `[]string` slice. The single dot `{{.}}` represents the string item in the loop.
* `{{range .Options}}<a href="/{{.Arc}}">{{.Text}}</a>{{end}}` – Loops over the `[]Option` slice, accessing nested fields `Arc` and `Text`.
* **Automatic XSS Prevention:** `html/template` automatically sanitizes values, protecting output against HTML/JavaScript injection attacks.

---

## 5. CLI Execution & Working Directories

* **Execution Context:** Paths passed to functions like `os.Open("gopher.json")` resolve relative to where `go run` is executed in the terminal—not where `main.go` resides on disk.
* **Flag Package (`flag`):** Use command-line flags to avoid hardcoded file paths and ensure binaries work cleanly across directory structures:

```go
filename := flag.String("file", "gopher.json", "path to the CYOA JSON file")
flag.Parse()
file, err := os.Open(*filename)

```

---