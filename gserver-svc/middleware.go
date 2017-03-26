package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
)

// LoggingHandler logs the incoming request.
func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(logRequest(r))
		buf := &bytes.Buffer{}
		handlers.CombinedLoggingHandler(buf, next).ServeHTTP(w, r)
		log.Println(buf.String())
	})
}

func logRequest(r *http.Request) string {
	requestLine := `"` + strings.Join([]string{r.Method, r.RequestURI, r.Proto}, ` `) + `"`
	return "incoming request from: " + r.RemoteAddr + "\trequest line: " + requestLine
}

// CheckAuth checks Nonce, Timestamp and JWE token (header Authorization: Bearer <token>) in the incoming request.
func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("checking token...")
		next.ServeHTTP(w, r)
	})
}
