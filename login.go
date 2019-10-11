package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

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
	r.ParseForm()

	Identity := tekmor.Identity{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	Details, err := Identity.Authenticate()

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}

	_, token, err := ta.Encode(jwt.MapClaims{"username": Details.Username, "home": Details.Home, "group": Details.Group})

	Cewkie := http.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	}

	http.SetCookie(w, &Cewkie)

	http.Redirect(w, r, "/access", http.StatusTemporaryRedirect)
}

func proccessJSONLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	type res struct {
		Token string `json:"token"`
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

	_, token, _ := ta.Encode(jwt.MapClaims{"username": Details.Username, "home": Details.Home, "group": Details.Group})

	Response := res{
		Token: token,
	}

	serialized, err := json.Marshal(Response)

	if err != nil {
		http.Error(w, "Unable to serialize content for response", http.StatusInternalServerError)
	}

	w.Write(serialized)

}

func displayLoginController(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Target string
	}

	t := template.New("login")
	t, err := t.Parse(loginPage)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}

	t.Execute(w, Response{
		Target: rootURL + "/login",
	})
}
