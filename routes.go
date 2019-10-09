package main

import (
	"net/http"
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

	r.Route("/login", func(r chi.Router) {

		r.Post("/", loginProcessorController)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(ta))

		//Covers both JSON and Non requests :3
		r.Use(combinedJWTMiddleware)

		r.Get("/access", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Well hello"))
		})
	})

	return r
}
