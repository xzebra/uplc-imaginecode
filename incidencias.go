package main

import (
	"time"
	"context"
	"net/http"
	"database/sql"
	"strconv"
)

type Incidencia struct {
	ID int
	Author string
	Type int
	Desc string
	Ubicacion string
	Estado int
	Apoyo string
	Date string
}

var Incidencias []Incidencia

func (i Incidencia) Insertar() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := DB.ExecContext(ctx,
		"INSERT INTO incidencias (id, autor, tipo, descripcion, ubicacion, estado, apoyo, timestamp) VALUES (?,?,?,?,?,?,?,?)",
		nil, i.Author, i.Type, i.Desc, i.Ubicacion, i.Estado, i.Apoyo, i.Date)
	if err != nil { return err }

	Incidencias = append(Incidencias, i)

	return nil
}

func ReadIncidencias() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cursor *sql.Rows
	cursor, err = DB.QueryContext(ctx,
		"SELECT id, autor, tipo, descripcion, ubicacion, estado, apoyo, timestamp FROM incidencias WHERE estado != 0")
	if err != nil { return err }
	defer cursor.Close()

	for cursor.Next() {
		var i Incidencia
		err = cursor.Scan(&i.ID, &i.Author, &i.Type, &i.Desc, &i.Ubicacion, &i.Estado, &i.Apoyo, &i.Date)
		if err != nil { return err }

		Incidencias = append(Incidencias, i)
	}

	return nil
}

func addIncidencia(w http.ResponseWriter, r *http.Request) {
	typ, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		SendError(w, err.Error())
		return
	}

	estado, err := strconv.Atoi(r.FormValue("estado"))
	if err != nil {
		SendError(w, err.Error())
		return
	}

	incid := &Incidencia{
		Author: r.FormValue("author"), Type: typ,
		Desc: r.FormValue("desc"), Ubicacion: r.FormValue("ubicacion"),
		Estado:estado, Date: time.Now().String(),
	}
	err = incid.Insertar()
	if err != nil {
		SendError(w, err.Error())
	} else {
		w.Write([]byte("ok"))
	}
}

func listarIncidencias(w http.ResponseWriter, r *http.Request) {
	ServeTemplate(w, "listado.html", struct {
		Lista []Incidencia
	}{
		Lista: Incidencias,
	})
}

func handleIncidencias(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listarIncidencias(w, r)
	case "POST":
		addIncidencia(w, r)
	}
}
