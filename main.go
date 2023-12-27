package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var cfg = Config{
	requestCounter: 0,
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	router := chi.NewRouter()

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

	router.Get("/healthz", readinessHandler)

	router.Get("/metrics", cfg.printRequestsHandler)

	router.HandleFunc("/reset", cfg.resetHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: routerCors,
	}

	// corsMux.ServeHTTP(http.ResponseWriter, *http.Request)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
