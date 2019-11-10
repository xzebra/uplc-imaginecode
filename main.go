package main

import (
	"os"
	"fmt"
	"net/http"
	"text/template"
	"path/filepath"
)

var (
	demoKey = "unipalaco"
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

func SendError(w http.ResponseWriter, err string) {
	ServeTemplate(w, "error.html", struct {
		Error string
	}{
		err,
	})
}

func comprobarAcceso(w http.ResponseWriter, r *http.Request) bool {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "authKey" {
			return cookie.Value == demoKey
		}
	}
	SendError(w, "Debes estar <a href=\"/login.html\">loggeado</a> para acceder")
	return false
}

func main() {
	DBInit()
	defer DB.Close()

	// pre guardar las incidencias y
	// actualizar en un futuro
	err := ReadIncidencias()
	if err != nil {
		panic(err)
	}

	fmt.Println("running server")

	http.HandleFunc("/crear_incidencia.html", handleCrearIncidencia)

	http.HandleFunc("/incidencias", handleIncidencias)
	http.HandleFunc("/map", handleMap)

	http.HandleFunc("/admin", handleAdmin)
	http.HandleFunc("/admin/incidencias", handleAdminIncidencias)
	http.HandleFunc("/admin/map", handleAdminMap)
	http.HandleFunc("/admin/users", handleAdminUsers)

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
