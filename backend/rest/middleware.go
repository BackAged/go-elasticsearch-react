package rest

import (
	"net/http"

	"github.com/BackAged/go-elasticsearch-react/backend/config"
)

// APIKeyOnly protects apis by api key
func APIKeyOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cnf := config.GetApp()
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey != cnf.APIKey {
			ServeJSON(w, "UnAuthorized", http.StatusForbidden, "missing auth token", nil, nil, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
