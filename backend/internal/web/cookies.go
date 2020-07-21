package web

import (
	"net/http"
	"strings"
)

// AddCookie adds a new session cookie
func (app App) AddCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     app.ckProps.Name,
		Value:    sessionID,
		HttpOnly: app.ckProps.HTTPOnly,
		Secure:   app.ckProps.Secure,
	}
	ss := strings.ToLower(app.ckProps.SameSite)
	switch ss {
	case "lax":
		cookie.SameSite = http.SameSiteLaxMode
	case "strict":
		cookie.SameSite = http.SameSiteStrictMode
	default:
		cookie.SameSite = http.SameSiteDefaultMode
	}
	http.SetCookie(w, cookie)
	return
}

// ClearCookie sets a cookie to be deleted
func (app App) ClearCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     app.ckProps.Name,
		Value:    "",
		HttpOnly: app.ckProps.HTTPOnly,
		Secure:   app.ckProps.Secure,
		MaxAge:   0,
	}
	http.SetCookie(w, cookie)
	return
}
