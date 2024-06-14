package main

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey int

const (
	contextKeyUserName contextKey = iota
)

func (app *application) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func (app *application) authRequiredAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Verify if user has the right to access the resource
		if !claims.HasRight("admin") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Adding username within context
		ctx := context.WithValue(r.Context(), "username", claims.Name)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) authRequiredProtected(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Verify if user has the right to access the resource
		if !claims.HasRight("user") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// Adding username within context
		ctx := context.WithValue(r.Context(), contextKeyUserName, claims.Name)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
