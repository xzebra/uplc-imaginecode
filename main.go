package main

import (
	"net/http"
)

func handle(w ResponseWriter, r *Request) {

}

func main() {
	http.HandleFunc("/", handle)
	fmt.Println("running server")
	http.ListenAndServe(":8080", nil)
}
