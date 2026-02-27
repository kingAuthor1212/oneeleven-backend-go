package main

import (
	"log"
	"net/http"
	"oneeleven-backend-go/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sort", CorsMiddleware(handlers.SortHandler))

	log.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}