package main

import (
	"fmt"
	"net/http"
	"text/template"
	"path/filepath"
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

func main() {
	DBInit()
	defer DB.Close()

	fmt.Println("running server")
	//http.HandleFunc("/incidencias", handleIncidencias)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe("localhost:8080", nil)
}
