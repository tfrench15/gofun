package urlshort

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	seen := false
	redir := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		url, ok := pathsToUrls[path]
		if ok {
			seen = true
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}
	}
	if seen {
		return redir
	}
	return func(w http.ResponseWriter, r *http.Request) {
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	ymlMap := make(map[string]string)
	err := yaml.Unmarshal(yml, &ymlMap)
	if err != nil {
		log.Fatalf("Error reading YAML: %v", err)
		return nil, err
	}
	fmt.Println(ymlMap)
	return nil, nil
}
