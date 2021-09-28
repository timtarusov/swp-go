package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ts.tarusov/swp/internal/handlers"
)

func loggingMiddleWare(srv *handlers.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			prid := r.URL.Query().Get("project_id")
			uri := r.RequestURI
			method := r.Method
			srv.Info("Handling: \n\tURI: %s\n\tMethod: %s\n\tProjectId: %s\n", uri, method, prid)

			next.ServeHTTP(w, r)
		})
	}
}

func authorizationMiddleWare(srv *handlers.Service) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			api_token := os.Getenv("API_SECRET")
			if header != fmt.Sprintf("Bearer %s", api_token) {
				w.WriteHeader(http.StatusForbidden)
				if err := json.NewEncoder(w).Encode("Missing auth token"); err != nil {
					srv.Error("%v", err)
				}
				srv.Info("Authentication failed: %s", r.RequestURI)
			} else {
				srv.Info("Successfully authenticated")
			}
			next.ServeHTTP(w, r)
		})
	}
}
