package main

import (
	"fmt"
	"log"
	"net/http"

	urlshortner "gophercises/02-url-shortner"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func fallback(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "fallback, url not found in map\n")
}

func main() {
	mux := http.NewServeMux()

	// Register specific endpoints
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/headers", headers) 

	// Catch-all route for paths not found in mux
	mux.HandleFunc("/", fallback)

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/github":         "https://github.com",
		"/gophercises":    "https://gophercises.com",
		"/go":             "https://golang.org",
	}

	mapHandler := urlshortner.MapHandler(pathsToUrls, mux)

	yml := `
- path: /urlshort-godoc-yaml
  url: https://godoc.org/github.com/gophercises/urlshort
- path: /yaml-godoc-yaml
  url: https://godoc.org/gopkg.in/yaml.v2
`

	yamlHandler, err := urlshortner.YAMLHandler([]byte(yml), mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8090...")
	http.ListenAndServe(":8090", yamlHandler)
}
