package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/jwtauth"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var ta *jwtauth.JWTAuth
var port int = 9666
var listenDirective string
var directory Directions
var directoryFile string = "directory.yml"
var rootURL string = ""

func main() {
	setup()
	r := routes()
	http.ListenAndServe(listenDirective, r)
}

func setup() {
	ta = jwtauth.New("HS256", []byte(os.Getenv("TOKEN_SECRET")), nil)

	if os.Getenv("LISTEN_PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("LISTEN_PORT"), 10, 64)
		if err != nil {
			panic("An unsuitable LISTEN_PORT was provided. Cannot setup")
		}

		port = int(p)
	}

	//Look for configured variable for loading the directory file
	if os.Getenv("DIRECTORY_FILE") != "" {
		directoryFile = os.Getenv("DIRECTORY_FILE")
	}

	buildDirectory()

	//Look for a custom ROOT URL
	//Defaults to "" so that when appended to /login you'll just get /login
	//But if a custom root of "/protected" is provided, login will redirect to "/protected/login"
	if os.Getenv("FOCAL_ROOT") != "" {
		rootURL = os.Getenv("FOCAL_ROOT")
	}

	listenDirective = ":" + strconv.Itoa(port)
	log.Print(listenDirective)
}

func buildDirectory() {
	if _, err := os.Stat(directoryFile); err == nil {
		log.Info("Located a directory file to parse")

		contents, err := ioutil.ReadFile(directoryFile)
		if err != nil {
			log.Error(err)
			directory = Directions{}
			return
		}

		directory = Directions{}

		err = yaml.Unmarshal(contents, &directory)
		if err != nil {
			panic("Unable to parse the listings! Giving up in a cowardly fashion")
		}
	}
}

func badAuthResponse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") == "application/json" {
		log.Error("No token is present and identified as a JSON request")
		log.Error("Request identified as ", r.Header.Get("content-type"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//For interactive Sessions, re-direct to / for login
	log.Info("No token present, but it appears to be an interactive session. Redirecting to / to login")
	http.Redirect(w, r, rootURL+"/", http.StatusTemporaryRedirect)
	return
}
