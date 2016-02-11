package main

import (
	"fmt"
	"log"
	"os"
	"pmoroney/poplop/db"

	"github.com/camlistore/camlistore/pkg/misc/pinentry"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error

	db.Connect()

	n, err := db.GetScheme(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	pe := pinentry.Request{
		Desc:   "Poplop",
		Prompt: "Master Passphrase",
		OK:     "OK",
		Cancel: "Cancel",
		Error:  "",
	}

	master, err := pe.GetPIN()
	if err != nil {
		log.Fatal(err)
	}

	pass, err := n.Hash(master)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pass)
}
