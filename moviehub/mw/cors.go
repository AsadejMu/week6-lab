
package mw

import (
	"net/http"
	"strings"
)

// CORS อนุญาตเฉพาะ origin ที่อยู่ในรายชื่อ
func CORS(allowed []string, next http.Handler) http.Handler {
	ok := map[string]bool{}

	for _, o := range allowed {
		o = strings.TrimSpace(o)
		if o != "" {
			ok[o] = true
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))

		if ok[origin] {
			h := w.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Add("Vary", "Origin")
			h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			h.Set("Access-Control-Max-Age", "600")
		}

		if r.Method == http.MethodOptions {
			if origin != "" && !ok[origin] {
				http.Error(w, "Origin not allowed", http.StatusForbidden)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
