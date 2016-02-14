package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pmoroney/poplop/db"

	"github.com/camlistore/camlistore/pkg/misc/pinentry"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

func main() {
	var err error

	if len(os.Args) != 2 {
		log.Fatal("Usage: poplop nickname\n")
	}

	cfg := struct {
		DB mysql.Config
	}{}

	cfgBytes, err := ioutil.ReadFile("/etc/poplop.conf")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	db.Connect(cfg.DB)

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
