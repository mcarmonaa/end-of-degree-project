package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mcarmonaa/end-of-degree-project/auth-svc/auth"
	"github.com/mcarmonaa/end-of-degree-project/message"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	// TODO: on deployment, service address must be taken from net.LookUpHost("servicename")
	authSvc    = os.Getenv("AUTH_SVC")
	usersSVC   = os.Getenv("USERS_SVC")
	devicesSVC = os.Getenv("DEVICES_SVC")
	// TODO: TLS certs?
	dialOptions = []grpc.DialOption{grpc.WithInsecure()}
)

// GetSalt handles incoming HTTP GET requests to /salt?mail=mail@domain. It returns the salt to get
// the derivated key for the mail's owner password.
func GetSalt(w http.ResponseWriter, r *http.Request) {
	mail := r.FormValue("mail")
	if mail == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(authSvc, dialOptions...)
	if err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}
	defer conn.Close()

	client := auth.NewAuthClient(conn)
	saltRequest := &auth.SaltRequest{Mail: mail}
	reply, err := client.GetSalt(context.Background(), saltRequest)
	if err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}

	salt := &message.AuthUser{Salt: reply.GetSalt()}
	response, err := json.Marshal(salt)
	if err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}

// Login handles incoming HTTP POST requests to /login. It authenticates a user in the system and returns a JWE token.
func Login(w http.ResponseWriter, r *http.Request) {
	loginMessage := &message.Encrypted{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(loginMessage); err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}

	conn, err := grpc.Dial(authSvc, dialOptions...)
	if err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}
	defer conn.Close()

	client := auth.NewAuthClient(conn)
	loginRequest := &auth.LoginRequest{
		Mail:    loginMessage.Mail,
		Iv:      loginMessage.IVector,
		Payload: loginMessage.Payload,
	}

	reply, err := client.Login(context.Background(), loginRequest)
	if err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}

	response := &message.Encrypted{
		IVector: reply.GetIv(),
		Payload: reply.GetPayload(),
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

// AddUser handles incoming HTTP POST requests to /register. It creates a new user in the system.
func AddUser(w http.ResponseWriter, r *http.Request) {
	newUser := &message.AuthUser{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(newUser); err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}

	conn, err := grpc.Dial(authSvc, dialOptions...)
	if err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}
	defer conn.Close()

	client := auth.NewAuthClient(conn)
	regRequest := &auth.RegisterRequest{Mail: newUser.Mail, Password: newUser.Password}
	reply, err := client.Register(context.Background(), regRequest)
	if err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}

	salt := &message.AuthUser{Salt: reply.GetSalt()}
	response, err := json.Marshal(salt)
	if err != nil {
		logAndSetHTTPErrorCode(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
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
