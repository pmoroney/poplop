package main

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/pmoroney/poplop/db"
)

func init() {
	http.HandleFunc("/poplop", poplop_handle)
}

func poplop_handle(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nick := req.Form.Get("nickname")
	master := req.Form.Get("master")
	if nick == "" || master == "" {
		http.Error(w, "Both nickname and master parameters are required", http.StatusBadRequest)
		return
	}

	n, err := db.GetScheme(nick)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err == sql.ErrNoRows {
		http.Error(w, "Nickname not found", http.StatusNotFound)
		return
	}

	pass, err := n.Hash(master)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, pass)
}
