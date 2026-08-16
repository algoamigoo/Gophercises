package urlshortner

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			fmt.Printf("Redirecting %s -> %s\n", path, dest)
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		// Path not found in map, delegate to fallback handler
		fmt.Printf("Path %q not found in map, forwarding to fallback\n", path)
		fallback.ServeHTTP(w, r)
	}
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathURLs)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYAML(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

func buildMap(pathURLs []pathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathURLs {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}
