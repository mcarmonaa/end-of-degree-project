package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/mcarmonaa/end-of-degree-project/auth-svc/auth"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/grpc"
)

const dbDialect = "postgres"

func init() {
	var loggingOut = os.Stdout
	log.SetPrefix("[ auth-svc ] ")
	log.SetOutput(loggingOut)
	log.SetFlags(log.LstdFlags | log.LUTC)
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":80", "IP:port to bind")
	flag.Parse()

	db, err := gorm.Open(dbDialect, os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	service, err := auth.NewAuthSvc(db)
	if err != nil {
		log.Fatal(err)
	}

	auth.RegisterAuthServer(server, service)
	log.Println("Listening...")
	log.Fatal(server.Serve(lis))
}
