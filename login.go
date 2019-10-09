package main

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/metrumresearchgroup/tekmor"
	log "github.com/sirupsen/logrus"
)

//Do the operation of logging in. Split actions based on content type provided.
func loginProcessorController(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		processFormLogin(w, r)
		return
	}

	proccessJSONLogin(w, r)
}

//The whole purpose is to provide a login mechanism for interactive avenues where a cookie will be used
//This adheres to the chi authjwt model interactively
func processFormLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This looks like a form login!"))
}

func proccessJSONLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	type res struct {
		Token string `json:"string"`
	}

	Identity := tekmor.Identity{}

	err := json.NewDecoder(r.Body).Decode(&Identity)
	if err != nil {
		http.Error(w, "Invalid request content. Could not be serialized", http.StatusBadRequest)
		return
	}

	//Attempt Shell Login
	Details, err := Identity.Authenticate()

	if err != nil {
		log.Error(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	_, token, err := ta.Encode(jwt.MapClaims{"username": Details.Username, "home": Details.Home, "group": Details.Group})

	w.Write([]byte(token))

}