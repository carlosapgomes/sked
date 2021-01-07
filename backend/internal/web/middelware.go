package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/text/language"
)

func (app App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// verify if there is a cookie with a session ID
		cookie, err := r.Cookie(app.ckProps.Name)
		if err != nil {
			// there is no cookie, just call next handlers
			next.ServeHTTP(w, r)
			return
		}
		sID := cookie.Value
		// verify if this session ID still exists in the database
		// and has not expired
		session, err := app.sessionService.Read(sID)
		if err != nil {
			// there is no session with this sID in the DB
			// remove current cookie and call next handler
			app.ClearCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		if time.Now().UTC().After(session.ExpiresAt) {
			// current session has expired
			// remove cookie and delete this session from DB
			app.ClearCookie(w)
			err := app.sessionService.Delete(session.ID)
			if err != nil {
				app.serverError(w, err)
			}
			next.ServeHTTP(w, r)
			return
		}

		// verify if the user associated with the current session is active
		user, err := app.userService.FindByID(session.UID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if !user.Active {
			// current user is not active anymore,
			// remove cookie and delete this session from DB
			app.ClearCookie(w)
			err := app.sessionService.Delete(session.ID)
			if err != nil {
				app.serverError(w, err)
			}
			next.ServeHTTP(w, r)
			return
		}
		// set context keys for authenticated user, sessionID and user.Roles
		ctx := context.WithValue(r.Context(), ContextKeyIsAuthenticated, true)
		ctx = context.WithValue(ctx, ContextKeyUser, user)
		ctx = context.WithValue(ctx, ContextKeySID, session.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app App) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			app.clientError(w, http.StatusForbidden)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (app App) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAdmin(r) {
			app.clientError(w, http.StatusForbidden)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// func (app App) requireSameUserOrAdmin(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !app.IsSameUserOrAdmin(r) {
// 			app.clientError(w, http.StatusForbidden)
// 			return
// 		}
// 		w.Header().Add("Cache-Control", "no-store")
// 		next.ServeHTTP(w, r)
// 	})
// }

// The bellow middlewares were
// shamelessly copied from "Let'sGo!" book by
// Alex Edwards (https://www.alexedwards.net/#book)
func (app App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app App) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app App) detectLang(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		// the default language will be selected for t == nil
		tag, _, _ := app.langMatcher.Match(t...)
		ctx := context.WithValue(r.Context(), ContextKeyLang, tag.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
