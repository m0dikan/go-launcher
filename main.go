package launcher

import (
	"fmt"
	"net/http"
)

func runmain() {
	fmt.Println("server started")

	mux := http.NewServeMux()
	mux.HandleFunc("/api", proxy(handler))
}

func proxy(f func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("request")
		f(w, req)
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}
