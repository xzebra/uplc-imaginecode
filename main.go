package main

import (
	"net/http"
)

func ServeTemplate(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filepath.Join("templates", filename))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func handleIncidencias(w ResponseWriter, r *Request) {
	ServeTemplate(w, "client_incidence"
}

func main() {
	fmt.Println("running server")
	http.HandleFunc("/incidencias", handleIncidencias)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":8080", nil)
}
