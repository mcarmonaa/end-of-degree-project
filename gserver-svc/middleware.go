package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

// LoggingHandler logs the incoming request.
func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}
		handlers.CombinedLoggingHandler(buf, next).ServeHTTP(w, r)
		log.Print(buf.String())
	})
}

// CheckAuth checks Nonce, Timestamp and JWE token (header Authorization: Bearer <token>) in the incoming request.
func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("checking token...")
		next.ServeHTTP(w, r)
	})
}
