package main

import (
	"fmt"
	"net/http"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Helo wolf</h2>")
	} else if r.URL.Path == "/dog" {
		fmt.Fprint(w, "<h1>Helo dog</h2>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>We could not find the page:(</h2>")
	}

}

func main() {
    http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":3000", nil)
}
