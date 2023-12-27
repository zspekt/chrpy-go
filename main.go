package main

import (
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	requestCounter int
}

var cfg = Config{
	requestCounter: 0,
}

func (c *Config) trackRequestWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.requestCounter += 1
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	mux.Handle(
		"/app/",
		http.StripPrefix("/app", cfg.trackRequestWrapper(http.FileServer(http.Dir(filepathRoot)))),
	)
	mux.HandleFunc("/healthz", readinessHandler)

	mux.HandleFunc("/metrics", cfg.printRequestsHandler)

	mux.HandleFunc("/reset", cfg.resetHandler)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	// corsMux.ServeHTTP(http.ResponseWriter, *http.Request)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (c *Config) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.requestCounter = 0
	num := strconv.Itoa(c.requestCounter)
	w.Write([]byte(num))
}

func (c *Config) printRequestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	num := strconv.Itoa(c.requestCounter)
	str := "Hits: " + num
	w.Write([]byte(str))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
