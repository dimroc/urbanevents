package main

import (
	//"encoding/json"
	"fmt"
	//"github.com/dimroc/urban-events/cityrecorder/cityrecorder"
	//"io/ioutil"
	//"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
