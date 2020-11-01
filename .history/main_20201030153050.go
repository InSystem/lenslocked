package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
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

// Hello is func that saying hello
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
    router.GET("/hello/:name", Hello)
	http.ListenAndServe(":3000", router)
}
