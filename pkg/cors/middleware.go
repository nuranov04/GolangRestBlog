package cors

import (
	"net/http"
)

func MiddleCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Max-Age", "15")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r)
	}
}
