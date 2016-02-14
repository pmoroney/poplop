package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pmoroney/poplop/db"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

//log.Printf("Request: %+v\n", req.TLS.VerifiedChains[0][0].Subject.CommonName)

//var templates = template.Must(template.ParseGlob("*.html"))

/*
func index(w http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(w, "index", struct{ Title string }{Title: "index"})
}
*/

func main() {
	cfg := struct {
		DB   mysql.Config
		HTTP struct {
			Addr              string
			ReadTimeout       time.Duration
			WriteTimeout      time.Duration
			KeyFile           string
			CertFile          string
			ClientCAFile      string
			RequireClientCert bool
		}
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

	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(cfg.HTTP.CertFile, cfg.HTTP.KeyFile)
	if err != nil {
		log.Fatal(err)
	}

	config.BuildNameToCertificate()

	if cfg.HTTP.RequireClientCert {
		config.ClientAuth = tls.RequireAndVerifyClientCert

		clientCAFile, err := ioutil.ReadFile(cfg.HTTP.ClientCAFile)
		if err != nil {
			log.Fatal(err)
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(clientCAFile)
		config.ClientCAs = pool
	}

	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		TLSConfig:    config,
	}

	log.Println("Listening on " + cfg.HTTP.Addr)
	server.ListenAndServeTLS("", "")
}
