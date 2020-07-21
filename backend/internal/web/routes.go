package web

import (
	"net/http"
)

// Routes function sets all mux routes
func (app *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/users/login", app.login()) // post
	// mux.Handle("/users/logout", app.requireAuthentication(app.logout())) // post
	mux.Handle("/users/validateEmail", app.validateEmail()) // ?token=xxxxxxxxx
	mux.Handle("/users/resetPassword", app.resetPassword()) // POST {"email": "email@address"}
	mux.Handle("/users/verifyResetPw", app.verifyResetPw()) //  ?token=xxxxxxxxx
	// GET "/users/$UID/name" // GET user's email address
	// GET "/users/$UID/email" // GET user's name
	// GET "/users/$UID" // GET user's data
	// TODO: those bellow two lines might not be neccessary
	// verify if "/users/" is enough to handle all "users" routes
	mux.Handle("/users/", app.requireAuthentication(app.users())) // ANY VERB
	mux.Handle("/users", app.requireAuthentication(app.users()))  // ANY VERB
	// for an open app
	// IMPORTANT: user signUp does not add a role
	// mux.Handle("/users/signup", app.signUp()) //post

	// for a more strict app, when only some users (i.e.: admins) can create a new user
	// IMPORTANT: user creation does not add a role
	// should be able to add new users:
	// mux.Handle("/users", app.requireAuthentication(app.requireAdmin(app.addUser()))) //post

	// health check route
	mux.Handle("/healthz", app.Healthz())
	return app.recoverPanic(app.logRequest(app.authenticate(mux)))
}
