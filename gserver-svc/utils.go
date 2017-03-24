package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func initRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	commonMid := alice.New(LoggingHandler)

	registerMid := commonMid.Append()
	router.Handle("/register", registerMid.ThenFunc(AddUser)).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	loginMid := commonMid.Append()
	router.Handle("/login", loginMid.ThenFunc(GetSalt)).Methods(http.MethodGet)
	router.Handle("/login", loginMid.ThenFunc(Login)).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	usersMid := commonMid.Append(CheckAuth)
	// users
	router.Handle("/users/{id:[0-9]+}", usersMid.ThenFunc(GetUser)).Methods(http.MethodGet)
	router.Handle("/users/{id:[0-9]+}", usersMid.ThenFunc(UpdateUser)).Methods(http.MethodPut).Headers("Content-Type", "application/json")
	router.Handle("/users/{id:[0-9]+}", usersMid.ThenFunc(RemoveUser)).Methods(http.MethodDelete)
	// devices associated to users
	router.Handle("/users/{id:[0-9]+}/devices", usersMid.ThenFunc(ListDevices)).Methods(http.MethodGet)
	router.Handle("/users/{id:[0-9]+}/devices", usersMid.ThenFunc(PairDevice)).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	router.Handle("/users/{id:[0-9]+}/devices/{id:[0-9]+}", usersMid.ThenFunc(GetDevice)).Methods(http.MethodGet)
	router.Handle("/users/{id:[0-9]+}/devices/{id:[0-9]+}", usersMid.ThenFunc(UnpairDevice)).Methods(http.MethodDelete)
	// mesasures from devices
	router.Handle("/users/{id:[0-9]+}/devices/{id:[0-9]+}/measures", usersMid.ThenFunc(ListMeasures)).Methods(http.MethodGet)
	router.Handle("/users/{id:[0-9]+}/devices/{id:[0-9]+}/measures", usersMid.ThenFunc(RemoveMeasures)).Methods(http.MethodDelete)

	devicesMid := commonMid.Append(CheckAuth)
	router.Handle("/devices/{id:[0-9]+}/measures", devicesMid.ThenFunc(UpdateMeasures)).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	return router
}

func debugRequest(r *http.Request) string {
	s := "\n"
	s += fmt.Sprintln(r.Method, r.RequestURI, r.Proto)
	s += fmt.Sprintln("Host:", r.Host)
	for k, v := range r.Header {
		s += fmt.Sprintln(k+":", strings.Join(v, "; "))
	}

	return s
}
