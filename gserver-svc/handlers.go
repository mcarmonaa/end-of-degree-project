package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	authSvc    = os.Getenv("AUTH_SVC")
	usersSVC   = os.Getenv("USERS_SVC")
	devicesSVC = os.Getenv("DEVICES_SVC")
)

// GetSalt handles incoming HTTP GET requests to /login?mail=mail@domain. It returns the salt to get
// the derivated key for the mail's owner password.
func GetSalt(w http.ResponseWriter, r *http.Request) {
	// /login?mail=mail@domain
	mail := r.FormValue("mail")
	if mail == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, mail)
}

// Login handles incoming HTTP POST requests to /login. It authenticates a user in the system and returns a JWE token.
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// AddUser handles incoming HTTP POST requests to /register. It creates a new user in the system.
func AddUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// GetUser handles incoming HTTP GET requests to /users/{id:[0-9]+}. It searchs for a user and returns its associated data.
func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// UpdateUser handles incoming HTTP PUT requests to /users/{id:[0-9]+}. It updates the user's associated data.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// RemoveUser handles incoming HTTP DELETE requests to /users/{id:[0-9]+}. It removes a user from the system.
func RemoveUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// ListDevices handles incoming HTTP GET requests to /users/{id:[0-9]+}/devices. It returns the devices' list assciated to a user.
func ListDevices(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// PairDevice handles incoming HTTP POST requests to /users/{id:[0-9]+}/devices. It associates a device to a user.
func PairDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// GetDevice handles incoming HTTP GET requests to /users/{id:[0-9]+}/devices/{id:[0-9]+}. It returns the device's data associated to a user.
func GetDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// UnpairDevice handles incoming HTTP DELETE requests to /users/{id:[0-9]+}/devices/{id:[0-9]+}. It removes an invalidates a device from the user's list of associated devices.
func UnpairDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// ListMeasures handles incoming HTTP GET requests to /users/{id:[0-9]+}/devices/{id:[0-9]+}/measures. It returns the list of measures associated to a device.
func ListMeasures(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// RemoveMeasures handles incoming HTTP DELETE requests to /users/{id:[0-9]+}/devices/{id:[0-9]+}/measures. It removes the list of measures associated to a device.
func RemoveMeasures(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}

// UpdateMeasures handles incoming HTTP PUT requests to /devices/{id:[0-9]+}/measures. It add measures to the device's list of measures.
func UpdateMeasures(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RequestURI)
}
