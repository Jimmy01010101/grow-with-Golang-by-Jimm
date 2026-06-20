package main

import (
	"log"
	"net/http"

	"go-login/handlers"
	"go-login/store"
)

func main() {
	memStore := store.NewMemStore()

	authHandler := handlers.NewAuthHandler(memStore)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("GET /profile", handlers.AuthMiddleware(authHandler.Profile))

	addr := ":1234"
	log.Printf("Server jalan di http://localhost%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}