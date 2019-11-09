package main

import (
	"time"
	"context"
	"net/http"
	"database/sql"
)

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	ServeTemplate(w, "admin.html", struct {
		Lista []Incidencia
	}{
		Lista: Incidencias,
	})
}

func handleAdminMap(w http.ResponseWriter, r *http.Request) {	
	ServeTemplate(w, "mapAll.html", struct {
		Lista []Incidencia
	}{
		Lista: Incidencias,
	})
}

type User struct {
	Username string
	Fecha string
	Mail string
	Incidencias string
}

func ReadUsers() (users []User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cursor *sql.Rows
	cursor, err = DB.QueryContext(ctx,
		"SELECT username, fecha, mail, incidencias FROM usuarios WHERE 1")
	if err != nil { return nil, err }
	defer cursor.Close()

	for cursor.Next() {
		var i User
		err = cursor.Scan(&i.Username, &i.Fecha, &i.Mail, &i.Incidencias)
		if err != nil { return nil, err }

		users = append(users, i)
	}

	return users, nil
}

func handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ReadUsers()
	if err != nil { return }
	
	ServeTemplate(w, "users.html", struct {
		Lista []User
	}{
		Lista: users,
	})
}

func handleAdminIncidencias(w http.ResponseWriter, r *http.Request) {
	ServeTemplate(w, "incidencias.html", struct {
		Lista []Incidencia
	}{
		Lista: Incidencias,
	})
}
