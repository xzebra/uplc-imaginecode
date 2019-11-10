package main

import (
	"time"
	"context"
	"net/http"
	"path/filepath"
	"database/sql"
)

type Incidencia struct {
	ID int
	Author string
	Type string
	Desc string
	Longitud string
	Latitud string
	Estado string
	Apoyo int
	Date string
	EsGrave bool
}

var Estados = []string{
	"En proceso",
}

var Tipos = map[string]string{
	"1": "Rescate Urgente",
	"2": "Mal estado",
	"3": "Incendio",
	"4": "Caida",
	"5": "Otros",
}

var EsGrave = map[string]bool{
	"Rescate Urgente": true,
	"Mal estado": false,
	"Incendio": true,
	"Caida": true,
	"Otros": false,
}


var Incidencias []Incidencia

func (i Incidencia) Insertar() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := DB.ExecContext(ctx,
		"INSERT INTO incidencias (id, autor, tipo, descripcion, longitud, latitud, estado, apoyo, timestamp) VALUES (?,?,?,?,?,?,?,?,?)",
		nil, i.Author, i.Type, i.Desc, i.Longitud, i.Latitud, i.Estado, i.Apoyo, i.Date)
	if err != nil { return err }

	i.EsGrave = EsGrave[i.Type]

	Incidencias = append(Incidencias, i)

	_, err = DB.ExecContext(ctx, "UPDATE usuarios SET incidencias = incidencias+1 WHERE username = ?", i.Author)
	return err
}

func EliminarIncidencia(id int) error {
	i := 0
	for ; i < len(Incidencias); i++ {
		if Incidencias[i].ID == id {
			break
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := DB.ExecContext(ctx, "DELETE FROM incidencias WHERE id = ?", Incidencias[i].ID)
	if err != nil { return err }

	_, err = DB.ExecContext(ctx, "UPDATE usuarios SET incidencias = incidencias-1 WHERE id = ?", Incidencias[i].Author)
	if err != nil { return err }

	copy(Incidencias[i:], Incidencias[i+1:])
	Incidencias = Incidencias[:len(Incidencias)-1]

	return nil
}

func ReadIncidencias() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cursor *sql.Rows
	cursor, err = DB.QueryContext(ctx,
		"SELECT id, autor, tipo, descripcion, longitud, latitud, estado, apoyo, timestamp FROM incidencias WHERE 1")
	if err != nil { return err }
	defer cursor.Close()

	for cursor.Next() {
		var i Incidencia
		err = cursor.Scan(&i.ID, &i.Author, &i.Type, &i.Desc, &i.Longitud, &i.Latitud, &i.Estado, &i.Apoyo, &i.Date)
		if err != nil { return err }

		i.EsGrave = EsGrave[i.Type]

		Incidencias = append(Incidencias, i)
	}

	return nil
}

func addIncidencia(w http.ResponseWriter, r *http.Request) {
	incid := &Incidencia{
		Author: r.FormValue("author"), Type: Tipos[r.FormValue("type")],
		Desc: r.FormValue("desc"),
		Longitud: r.FormValue("long"),
		Latitud: r.FormValue("lat"),
		Estado: Estados[0], Date: time.Now().Format("Jan 2 3:04pm"),
	}
	err := incid.Insertar()
	if err != nil {
		SendError(w, err.Error())
	} else {
		listarIncidencias(w, r)
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
	if !comprobarAcceso(w, r) {
		return
	}

	switch r.Method {
	case "GET":
		listarIncidencias(w, r)
	case "POST":
		addIncidencia(w, r)
	}
}

func handleMap(w http.ResponseWriter, r *http.Request) {
	if !comprobarAcceso(w, r) {
		return
	}

	ServeTemplate(w, "map.html", struct {
		Lista []Incidencia
	}{
		Lista: Incidencias,
	})
}

func handleCrearIncidencia(w http.ResponseWriter, r *http.Request) {
	if !comprobarAcceso(w, r) {
		return
	}

	http.ServeFile(w, r, filepath.Join("public", "crear_incidencia.html"))
}


