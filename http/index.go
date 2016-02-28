package main

import (
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/", index_handle)
}

func index_handle(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="utf-8">
	<title>Poplop</title>
	</head>
	<body>
	<form action="/poplop" method="post">
	<input name="nickname"/>
	<input type="password" name="master"/>
	<button type="submit">Get Password</button>
	</form>
	</body>
	</html>
	`)
}
