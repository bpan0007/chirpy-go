package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve static files from the ./assets directory
	fileServer := http.FileServer(http.Dir("./assets"))
	// Use http.StripPrefix to remove '/assets' prefix before passing to fileServer
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	mux.Handle("/app/*", http.StripPrefix("/app/*", http.FileServer(http.Dir("/"))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// func healthzHandler(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("This is the about page."))
	// }
	corsMux := middlewareCors(mux)
	err := http.ListenAndServe(":8080", corsMux)
	if err != nil {
		log.Fatal(err)
	}
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
