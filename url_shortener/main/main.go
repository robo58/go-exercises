package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"url_shortener"
)


func main(){
	filePath:=flag.String("path", "", "File path for data, JSON and YAML supported")
	flag.Parse()
	mux := defaultMux()

	pathsToUrls := map[string]string {
		"/google": "https://google.com",
		"/youtube": "https://youtube.com",
	}
	mapHandler := url_shortener.MapHandler(pathsToUrls,mux)

	var handler http.HandlerFunc
	var handleError error

	fileType := ""

	data, readError := readFileByName(filePath)
	if readError == nil {
		fileType = strings.SplitAfter(*filePath, ".")[1]
	}

	switch fileType {
		case "yaml":
			handler, handleError = url_shortener.YAMLHandler(data, mapHandler)
			if handleError != nil {
				panic(handleError)
			}
		case "json":
			handler, handleError = url_shortener.JSONHandler(data, mapHandler)
			if handleError != nil {
				panic(handleError)
			}
		default:
			handler = mapHandler
	}


	fmt.Println("Starting server on :8080")
	httpError := http.ListenAndServe(":8080", handler)
	if httpError != nil {
		panic(httpError)
	}
}

func defaultMux() *http.ServeMux{
	mux:=http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request){
	_, err := fmt.Fprintln(w, "Hello World.")
	if err != nil {
		panic(err)
	}
}

func readFileByName(filename *string) ([]byte, error){
	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		return nil,err
	}
	return data, nil
}