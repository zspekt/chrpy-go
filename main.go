package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// me when when i dont know how to handle these things
var cfg Config = Config{
	requestCounter: 0,
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	router := chi.NewRouter()
	routerAPI := chi.NewRouter()
	routerAdmin := chi.NewRouter()

	router.Mount("/api/", routerAPI)
	router.Mount("/admin/", routerAdmin)

	routerCors := middlewareCors(router)
	// router.Handle(pattern string, handler http.Handler)
	router.Handle(
		"/app",
		http.StripPrefix("/app", cfg.trackRequestWrapper(http.FileServer(http.Dir(filepathRoot)))),
	)

	router.Handle(
		"/app/*",
		http.StripPrefix("/app", cfg.trackRequestWrapper(http.FileServer(http.Dir(filepathRoot)))),
	)

	routerAPI.HandleFunc("/reset", cfg.resetHandler)

	routerAPI.Get("/healthz", readinessHandler)

	// chirp handlers
	routerAPI.Get("/chirps", chirpsGetHandler)
	routerAPI.Get("/chirps/*", chirpsGetByIDHandler)
	routerAPI.Post("/chirps", chirpsPostHandler)
	routerAPI.Delete("/chirps/*", chirpsDeleteHandler)

	// user handlers
	routerAPI.Post("/users", usersPostHandler)
	routerAPI.Post("/login", usersAuthHandler)
	routerAPI.Put("/users", usersEditHandler)

	routerAdmin.Get("/metrics", cfg.printRequestsHandler)

	routerAPI.Post("/refresh", refreshPostHandler)
	routerAPI.Post("/revoke", revokePostHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: routerCors,
	}

	// corsMux.ServeHTTP(http.ResponseWriter, *http.Request)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
