package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Config struct {
	requestCounter int
}

// increments counter when a request is received
func (c *Config) trackRequestWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.requestCounter += 1
		next.ServeHTTP(w, r)
	})
}

// resets the counter that tracks requests
func (c *Config) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.requestCounter = 0
	num := strconv.Itoa(c.requestCounter)
	w.Write([]byte(num))
}

// displays what the request counter is currently at
// func (c *Config) printRequestsHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	num := strconv.Itoa(c.requestCounter)
// 	str := "Hits: " + num
// 	w.Write([]byte(str))
// }

// displays what the request counter is currently at
func (c *Config) printRequestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	str := fmt.Sprintf(`
    <html>

    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>

    </html>
    `, c.requestCounter)
	w.Write([]byte(str))
}

// health check
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
