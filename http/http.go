package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pmoroney/poplop/db"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/acme/autocert"
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
			Host              string
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

	if cfg.HTTP.Host == "" && (cfg.HTTP.KeyFile == "" || cfg.HTTP.CertFile == "") {
		log.Fatal("Need host or keyfile and certfile configs")
	}

	if cfg.HTTP.Host == "" {
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(cfg.HTTP.CertFile, cfg.HTTP.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		config.BuildNameToCertificate()
	} else {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.HTTP.Host),
			Cache:      autocert.DirCache("."),
		}

		config.GetCertificate = m.GetCertificate
	}

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

	if cfg.HTTP.RequireClientCert {
		copy := *config
		copy.ClientAuth = tls.VerifyClientCertIfGiven
		config.GetConfigForClient = func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
			if !strings.HasSuffix(hello.ServerName, ".acme.invalid") {
				return nil, nil
			}

			return &copy, nil
		}
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
