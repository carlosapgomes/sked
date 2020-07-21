package web

import (
	"net/http"
)

// Healthz is a health check handler
func (app App) Healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
