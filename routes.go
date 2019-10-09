package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gopkg.in/yaml.v2"
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

	r.Post("/login", loginProcessorController)

	r.Get("/test", serializationController)

	//Iterate over config
	for _, v := range directory {
		r.Get("/"+strings.ToLower(v.Name), generateProxyHandler(v))
	}

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(ta))

		//Covers both JSON and Non requests :3
		r.Use(combinedJWTMiddleware)

		r.Get("/access", listingController)
		r.Post("/access", listingController)

	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	return r
}

func listingController(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Well hello"))
}

func serializationController(w http.ResponseWriter, r *http.Request) {
	Directions := Directions{
		{
			Name:   "google",
			Target: "http://www.google.com",
		},
		{
			Name:   "yahoo",
			Target: "http://www.yahoo.com",
		},
	}

	serialized, err := yaml.Marshal(Directions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(serialized)
}

func generateProxyHandler(d Direction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//	origin, _ := url.Parse(d.Target)

		// 	director := func(req *http.Request) {
		// 		req.Header.Add("X-Forwarded-Host", req.Host)
		// 		req.Header.Add("X-Origin-Host", origin.Host)
		// 		req.URL.Scheme = "http"
		// 		req.URL.Host = origin.Host
		// 	}

		// 	proxy := &httputil.ReverseProxy{Director: director}

		// 	proxy.ServeHTTP(w, r)
		// }
		proxy(d.Target, w, r)
	}
}

func proxy(target string, w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(u)

	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Host = u.Host
	r.URL.Path = "/"

	proxy.ServeHTTP(w, r)

}
