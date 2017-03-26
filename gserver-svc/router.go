package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func initRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	commonMid := alice.New(LoggingHandler)

	registerMid := commonMid.Append()
	router.Handle("/register", registerMid.ThenFunc(AddUser)).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	saltMid := commonMid.Append()
	router.Handle("/salt", saltMid.ThenFunc(GetSalt)).Methods(http.MethodGet)

	loginMid := commonMid.Append()
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
	// upload measures endpoint for devices
	router.Handle("/devices/{id:[0-9]+}/measures", devicesMid.ThenFunc(UpdateMeasures)).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	return router
}
