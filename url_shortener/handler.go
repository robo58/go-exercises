package url_shortener

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		path := r.URL.Path
		// dest is value of map, ok is boolean that checks if given key exists in map
		if dest,ok := pathsToUrls[path]; ok {
			// if we can match a path
			// redirect to it
			http.Redirect(w, r, dest, http.StatusFound)
		}
		//else fallback
		fallback.ServeHTTP(w,r)
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
	// parse yaml
	pathUrls, parseError := parseYaml(yml)
	if parseError!=nil{
		return nil, parseError
	}
	pathsToUrls := buildMap(pathUrls)
	// return map handler using map
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error){
	pathUrls, parseError := parseJson(json)
	if parseError!=nil{
		return nil, parseError
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	// dont need index so _, pu is instance of pathUrl
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYaml(data []byte)([]pathUrl, error){
	var pathUrls []pathUrl
	marshalError := yaml.Unmarshal(data, &pathUrls)
	if marshalError != nil {
		return nil, marshalError
	}
	return pathUrls, nil
}

func parseJson(data []byte)([]pathUrl, error){
	var pathUrls []pathUrl
	marshalError := json.Unmarshal(data, &pathUrls)
	if marshalError != nil {
		return nil, marshalError
	}
	return pathUrls, nil
}

type pathUrl struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}