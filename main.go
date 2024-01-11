package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/zspekt/chrpy-go/internal/database"
)

var (
	cfg Config = Config{
		requestCounter: 0,
	}
	DATAB     *database.DB
	IdCounter int
)

func init() {
	var errr error
	DATAB, errr = database.NewDB("./database.json")
	if errr != nil {
		log.Fatal(errr)
	}
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

	routerAPI.Get("/healthz", readinessHandler)
	routerAPI.HandleFunc("/reset", cfg.resetHandler)
	routerAPI.Post("/chirps", chirpsHandler)

	routerAdmin.Get("/metrics", cfg.printRequestsHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: routerCors,
	}

	// corsMux.ServeHTTP(http.ResponseWriter, *http.Request)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
