package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var ta *jwtauth.JWTAuth
var port int = 9666
var listenDirective string

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

	listenDirective = ":" + strconv.Itoa(port)
	log.Print(listenDirective)
}

//This middleware is used for interactive targets (IE Grafana) that are
//completely and blissfully unaware of potential tokens. Here, failures to login
//or the absence of a token in session trigger a redirect to /login

//For API requests, we'll use the jwtauth.Authenticator middleware.
func interactiveSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
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
