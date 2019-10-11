package main

import (
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	log "github.com/sirupsen/logrus"
)

func routes() chi.Router {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", displayLoginController)

	r.Get("/test", listingController)

	r.Post("/login", loginProcessorController)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(ta))

		//Covers both JSON and Non requests :3
		r.Use(combinedJWTMiddleware)

		r.Get("/access", listingController)
		r.Post("/access", listingController)

		//Iterate over config
		for _, v := range directory {
			r.HandleFunc("/"+strings.ToLower(v.Name)+"*", generateProxyHandler(v))
		}

	})

	return r
}

func listingController(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Directory Directions
		RootURL   string
	}

	req := Request{
		Directory: directory,
		RootURL:   rootURL,
	}

	t := template.New("listing")
	t, err := t.Parse(backendListing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, req)
}

func generateProxyHandler(d Direction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy(d, w, r)
	}
}

func proxy(d Direction, w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(d.Target)

	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Host = u.Host

	pieces := strings.Split(r.URL.Path, "/")
	curPath := pieces[len(pieces)-1]

	//Handle proxy passing and only pass non top level paths to the downstream
	if strings.ToLower(curPath) == strings.ToLower(d.Name) {
		//This it he top level of the magical rainbow road
		r.URL.Path = "/"
	} else {
		r.URL.Path = strings.Replace(r.URL.Path, "/"+strings.ToLower(d.Name), "", 1)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.ServeHTTP(w, r)

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
