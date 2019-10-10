package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Well hello"))
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
