package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var ta *jwtauth.JWTAuth
var port int = 9666
var listenDirective string
var directory Directions
var directoryFile string = "directory.yml"

func main() {
	setup()
	r := routes()
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	FileServer(r, "/", http.Dir(filesDir))

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

//This middleware is used for interactive targets (IE Grafana) that are
//completely and blissfully unaware of potential tokens. Here, failures to login
//or the absence of a token in session trigger a redirect to /login

//For API requests, we'll use the jwtauth.Authenticator middleware.
func combinedJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			log.Error("An error occurred getting the token from the context: ", err)
			badAuthResponse(w, r)
			return
		}

		if token == nil || !token.Valid {
			badAuthResponse(w, r)
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
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
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}
