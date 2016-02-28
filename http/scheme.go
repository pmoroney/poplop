package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path"

	"github.com/pmoroney/poplop"

	"github.com/pmoroney/poplop/db"
)

func init() {
	http.HandleFunc("/scheme/", scheme_handle)
}

func scheme_handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		scheme_get(w, req)
	case "POST":
		scheme_post(w, req)
	case "PUT":
		scheme_put(w, req)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func scheme_get(w http.ResponseWriter, req *http.Request) {
	dir, name := path.Split(req.URL.Path)
	if dir != "/scheme/" {
		http.Error(w, "Unknown directory structure", http.StatusBadRequest)
		return
	}

	n, err := db.GetScheme(name)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err == sql.ErrNoRows {
		http.Error(w, "name not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func scheme_post(w http.ResponseWriter, req *http.Request) {
	dir, name := path.Split(req.URL.Path)
	if dir != "/scheme/" {
		http.Error(w, "Unknown directory structure", http.StatusBadRequest)
		return
	}
	if name == "" {
		http.Error(w, "name required /scheme/{name}", http.StatusBadRequest)
		return
	}

	n := poplop.Scheme{}
	err := json.NewDecoder(req.Body).Decode(&n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if n.Name == name {
		http.Error(w, "name must match url", http.StatusBadRequest)
	}

	err = db.InsertScheme(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
}

func scheme_put(w http.ResponseWriter, req *http.Request) {
	dir, name := path.Split(req.URL.Path)
	if dir != "/scheme/" {
		http.Error(w, "Unknown directory structure", http.StatusBadRequest)
		return
	}
	if name == "" {
		http.Error(w, "name required /scheme/{name}", http.StatusBadRequest)
		return
	}

	n := poplop.Scheme{}
	err := json.NewDecoder(req.Body).Decode(&n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateScheme(n, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}
