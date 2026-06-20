package handlers

import (
	"context"
	"net/http"
	"strings"

	"Go-golang/auth"
)

type contextKey string

const userContextKey contextKey = "username"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token tidak ada"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		username, err := auth.ParseToken(tokenStr)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token tidak valid"})
			return
		}
		
		ctx := context.WithValue(r.Context(), userContextKey, username)
		next(w, r.WithContext(ctx))
	}
}