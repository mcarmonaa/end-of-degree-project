package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	var loggingOut = os.Stdout
	log.SetPrefix("[ gserver ] ")
	log.SetOutput(loggingOut)
	log.SetFlags(log.LstdFlags | log.LUTC)
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":80", "IP:port to bind")
	flag.Parse()

	router := initRouter()
	server := &http.Server{
		// Addr: , if not set, port :80 in al interfaces by default
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	//log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	log.Print("Listening...")
	log.Fatal(server.ListenAndServe())
}
