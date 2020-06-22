package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

var HelloString string

func main() {
	var isSet bool
	HelloString, isSet = os.LookupEnv("HELLOSTRING")
	if isSet == false {
		HelloString = "Default hello!"
	}
	http.HandleFunc("/", HelloWorldHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, HelloString)
}