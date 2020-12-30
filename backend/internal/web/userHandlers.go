package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"carlosapgomes.com/sked/internal/session"
	"carlosapgomes.com/sked/internal/token"
	"carlosapgomes.com/sked/internal/user"
)

type userData struct {
	ID       string   `json:"ID,omitempty"`
	Name     string   `json:"Name"`
	Email    string   `json:"Email"`
	Phone    string   `json:"Phone"`
	Password string   `json:"Password"`
	Roles    []string `json:"Roles,omitempty"`
}

// validates request user data
func (u *userData) validate() url.Values {
	errs := url.Values{}

	// check if name empty
	if u.Name == "" {
		errs.Add("Name", "This field cannot be blank")
	}

	// check if password is empty
	// if u.Password == "" {
	// 	errs.Add("Password", "This field cannot be blank")
	// }
	// check the name field is between 3 to 120 chars
	nameLen := utf8.RuneCountInString(u.Name)
	if nameLen < 3 || nameLen > 100 {
		errs.Add("name", "The name field must be between 3-100 chars!")
	}

	// check Password min size
	// if utf8.RuneCountInString(u.Password) < 8 {
	// 	errs.Add("Password", "This field is too short (minimum is 8 characters)")
	// }
	// check if email is empty
	if u.Email == "" {
		errs.Add("Email", "This field cannot be blank")
	}

	if !validEmail(u.Email) {
		errs.Add("Email", "This field is invalid")
	}

	return errs
}

// validEmail returns true if email string has a valid length and format
func validEmail(email string) bool {
	// check for valid email based on https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return (len(email) <= 254 && rxEmail.MatchString(email))
}

// getUsers search for users based on query params
func (app App) getUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()
	Loop:
		for q := range v {
			switch {
			case strings.Contains(q, "email"):
				app.getUserByEmail(w, r)
				break Loop
			case strings.Contains(q, "name"):
				app.getUserByName(w, r)
				break Loop
			case strings.Contains(q, "previous"):
				app.getAllUsers(w, r)
				break Loop
			case strings.Contains(q, "next"):
				app.getAllUsers(w, r)
				break Loop
			default:
				app.clientError(w, http.StatusBadRequest)
			}
		}
	})

}

// users is the root handler for restful users endpoints
// A URL represents a parsed URL:
// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
func (app App) users() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case (path == "/users/logout"):
			app.logout().ServeHTTP(w, r)
		case strings.HasSuffix(path, "/users"):
			app.usersNoPath(w, r)
		case strings.HasSuffix(path, "/name"):
			app.userName(w, r)
		case strings.HasSuffix(path, "/email"):
			app.userEmail(w, r)
		case strings.HasSuffix(path, "/password"):
			app.userPassword(w, r)
		case strings.HasSuffix(path, "/doctors"):
			app.getAllDoctors(w, r)
		default:
			app.clientError(w, http.StatusBadRequest)
		}
	})
}

// getDoctors return a list of all doctors
func (app App) getAllDoctors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
	}
	docs, err := app.userService.GetAllDoctors()
	if err != nil {
		app.serverError(w, err)
		return
	}
	output, err := json.Marshal(docs)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// userPassword sets a user password (POST)
func (app App) userPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
	}
	uID := app.between(r.URL.Path, "/users/", "/password")
	if uID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if !app.isValidUUID(uID) {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// get logged in user from ctx
	u, ok := r.Context().Value(ContextKeyUser).(*user.User)
	if !ok {
		// no ContextKeyUser -> user is not authenticated
		app.clientError(w, http.StatusForbidden)
		return
	}
	if !(uID == u.ID) && !(app.findRole(u.Roles, user.RoleAdmin)) {
		// requesting user is neither the same user nor admin
		app.clientError(w, http.StatusForbidden)
		return
	}
	// get pw from body
	r.ParseForm()

	//b, err := ioutil.ReadAll(r.Body)
	//defer r.Body.Close()
	//if err != nil {
	//app.clientError(w, http.StatusBadRequest)
	//return
	//}
	type pw struct {
		Password string `json:"password"`
		Confirm  string `json:"confirm_password"`
	}
	var p pw
	p.Password = r.FormValue("password")
	p.Confirm = r.FormValue("confirm_password")
	//err = json.Unmarshal(b, &p)
	//if err != nil {
	//app.serverError(w, err)
	//return
	//}
	if (p.Password != p.Confirm) ||
		(p.Password == "") ||
		(utf8.RuneCountInString(p.Password) < 8) {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// invalidate current session
	sid, ok := r.Context().Value(ContextKeySID).(string)
	if !ok {
		// no ContextKeySID -> user is not authenticated
		app.clientError(w, http.StatusForbidden)
		return
	}
	app.ClearCookie(w)
	err := app.sessionService.Delete(sid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// update password
	err = app.userService.UpdatePw(uID, p.Password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// redirect user to operation-success page
	tplData := &templateData{
		Title: "Password updated",
		User: &user.User{
			Name:  u.Name,
			Email: u.Email,
			ID:    u.ID,
		},
		Link: "",
	}
	out, err := app.render("operation-success.gohtml", tplData)
	if err != nil {
		app.serverError(w, err)
		return
	}
	out.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// usersNoPath reroute based on verbs and queries
func (app App) usersNoPath(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.requireAdmin(app.getUsers()).ServeHTTP(w, r)
	case http.MethodPost:
		app.requireAdmin(app.addUser()).ServeHTTP(w, r)
	case http.MethodPut:
		app.requireAdmin(app.updateUser()).ServeHTTP(w, r)
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

// userEmail get user email
func (app App) userEmail(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.getUserEmail(w, r)
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app App) getUserEmail(w http.ResponseWriter, r *http.Request) {
	uID := app.between(r.URL.Path, "/users/", "/email")
	if uID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	ctx := context.WithValue(r.Context(), ContextResourceID, uID)

	if !app.IsSameUserOrAdmin(r.WithContext(ctx)) {
		app.clientError(w, http.StatusForbidden)
		return
	}
	u, err := app.userService.FindByID(uID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(&map[string]string{"email": u.Email})
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// userName set/change/read user name
func (app App) userName(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.getUserName(w, r)
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app App) getUserName(w http.ResponseWriter, r *http.Request) {
	uID := app.between(r.URL.Path, "/users/", "/name")
	if uID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// identify resource to be read
	ctx := context.WithValue(r.Context(), ContextResourceID, uID)

	// check authorization for this resource
	if !app.IsSameUserOrAdmin(r.WithContext(ctx)) {
		app.clientError(w, http.StatusForbidden)
		return
	}
	u, err := app.userService.FindByID(uID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(&map[string]string{"name": u.Name})
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// updateUser only updates users' phone and roles for now
func (app App) updateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		var updatedUser userData
		err = json.Unmarshal(b, &updatedUser)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// update user phone
		err = app.userService.UpdatePhone(updatedUser.ID, updatedUser.Phone)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// update user roles
		err = app.userService.UpdateRoles(updatedUser.ID, updatedUser.Roles)
		if err != nil {
			app.serverError(w, err)
			return
		}
	})
}

// in case an "/adduser" endpoint is needed
func (app App) addUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		var newUser userData
		err = json.Unmarshal(b, &newUser)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// do not validate for empty passwords
		if validationErrors := newUser.validate(); len(validationErrors) > 0 {
			err := map[string]interface{}{"validationError": validationErrors}
			w.Header().Set("Content-type", "application/json")
			app.clientError(w, http.StatusBadRequest)
			js, e := json.Marshal(err)
			if e != nil {
				app.serverError(w, e)
				return
			}
			w.Write(js)
			return
		}

		var uid *string
		uid, err = app.userService.Create(
			newUser.Name,
			newUser.Email,
			newUser.Password,
			newUser.Phone,
			newUser.Roles)
		if err != nil {
			if errors.As(err, &user.ErrDuplicateField) {
				w.Header().Set("Content-type", "application/json")
				app.clientError(w, http.StatusBadRequest)
				js, err := json.Marshal(&map[string]string{"email": "address already in use"})
				if err != nil {
					app.serverError(w, err)
					return
				}
				w.Write(js)
				return
			}
			app.serverError(w, err)
			return
		}

		// send msg/token to confirm user's email account
		// create a new ValidateEmail token
		token, err := app.tokenService.New(*uid, token.ValidateEmail)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// compose email
		link := fmt.Sprintf("https://%s%s%s",
			r.Host, "/api/users/validateEmail?token=", *token)
		tplData := &templateData{
			Title: "Email Validation",
			User: &user.User{
				Name: newUser.Name,
			},
			Link: link,
		}
		out, err := app.render("validate-email.gohtml", tplData)
		if err != nil {
			app.serverError(w, err)
			return
		}
		msg := out.String()
		// send to newUser.Email
		s := "Validação de Email"
		res, err := app.mailerService.Send(newUser.Name, newUser.Email, s, msg)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.infoLog.Printf("Sent %v to %v\nResponse Code: %v\nMsg: %v\n",
			s, newUser.Email, res.Code, res.Msg)
		newUser.Password = ""
		newUser.ID = *uid
		output, err := json.Marshal(newUser)
		if err != nil {
			app.serverError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	})
}

// signUp method (POST)
// func (app App) signUp() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}
// 		// Read body
// 		b, err := ioutil.ReadAll(r.Body)
// 		defer r.Body.Close()
// 		if err != nil {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}
// 		var newUser userData
// 		err = json.Unmarshal(b, &newUser)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		// check other fields
// 		validationErrors := newUser.validate()
// 		// check if password is empty
// 		if newUser.Password == "" {
// 			validationErrors.Add("Password", "This field cannot be blank")
// 		}
// 		if len(validationErrors) > 0 {
// 			err := map[string]interface{}{"validationError": validationErrors}
// 			w.Header().Set("Content-type", "application/json")
// 			app.clientError(w, http.StatusBadRequest)
// 			js, e := json.Marshal(err)
// 			if e != nil {
// 				app.serverError(w, e)
// 				return
// 			}
// 			w.Write(js)
// 			return
// 		}

// 		var uid *string

// 		uid, err = app.userService.Create(newUser.Name, newUser.Email, newUser.Password, newUser.Phone)
// 		if err != nil {
// 			if errors.As(err, &user.ErrDuplicateField) {
// 				w.Header().Set("Content-type", "application/json")
// 				app.clientError(w, http.StatusBadRequest)
// 				js, err := json.Marshal(&map[string]string{"email": "address already in use"})
// 				if err != nil {
// 					app.serverError(w, err)
// 					return
// 				}
// 				w.Write(js)
// 				return
// 			}
// 			app.serverError(w, err)
// 			return
// 		}
// 		newUser.Password = ""
// 		newUser.ID = *uid
// 		output, err := json.Marshal(newUser)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return
// 		}
// 		w.Write(output)
// 	})
// }

// Login method (POST)
func (app App) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		var data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := decodeJSONBody(w, r, &data)
		if err != nil {
			var mr *malformedRequest
			if errors.As(err, &mr) {
				http.Error(w, mr.msg, mr.status)
			} else {
				app.serverError(w, err)
			}
			return
		}
		// Never allow a user to log in with an empty password.
		// This will break security if the account validation flow
		// is interrupted before the user has set the password.
		// Because at this point the user will have an account
		// activated with an empty password.
		if !validEmail(data.Email) || (data.Password == "") {
			err := map[string]interface{}{"validationError": "invalid email and/or password"}
			w.Header().Set("Content-type", "application/json")
			app.clientError(w, http.StatusBadRequest)
			js, e := json.Marshal(err)
			if e != nil {
				app.serverError(w, e)
				return
			}
			w.Write(js)
			return
		}
		// call user service authenticate
		uid, err := app.userService.Authenticate(data.Email, data.Password)
		if err != nil {
			if err == user.ErrInvalidCredentials {
				app.clientError(w, http.StatusUnauthorized)
			} else {
				app.serverError(w, err)
			}
			return
		}

		user, err := app.userService.FindByID(*uid)
		if err != nil {
			app.serverError(w, err)
			return
		}
		if !user.Active {
			err := map[string]interface{}{"validationError": "User is not active"}
			w.Header().Set("Content-type", "application/json")
			app.clientError(w, http.StatusUnauthorized)
			js, e := json.Marshal(err)
			if e != nil {
				app.serverError(w, e)
				return
			}
			w.Write(js)
			return
		}
		if !user.EmailWasValidated {
			err := map[string]interface{}{"validationError": "Email was not validated"}
			w.Header().Set("Content-type", "application/json")
			app.clientError(w, http.StatusUnauthorized)
			js, e := json.Marshal(err)
			if e != nil {
				app.serverError(w, e)
				return
			}
			w.Write(js)
			return
		}

		sid, err := app.sessionService.Create(*uid)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.AddCookie(w, *sid)

		// do not return the password
		user.HashedPw = []byte("")
		js, err := json.Marshal(user)
		if err != nil {
			app.serverError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
}

// Logout Logs a user out (POST)
func (app App) logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		// find current session, returns error if not present
		c, err := r.Cookie(app.ckProps.Name)
		if err != nil {
			// because we use requireAuthentication in routes, http will
			// never reach this point without a cookie.
			// So this fragment never gets covered by tests
			if err == http.ErrNoCookie {
				msg := map[string]string{"Msg": "User is not logged in"}
				app.sendMsg(w, &msg)
			} else {
				app.serverError(w, err)
			}
			return
		}
		sID := c.Value
		// delete current session
		err = app.sessionService.Delete(sID)
		if err != nil && err != session.ErrNoRecord {
			app.serverError(w, err)
			return
		}
		// clear cookie
		app.ClearCookie(w)
		msg := map[string]string{"Msg": "User logged out"}
		app.sendMsg(w, &msg)
		return
	})
}

// validateEmail validate user's account email
func (app App) validateEmail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		// extract token from query
		tokenID := r.URL.Query().Get("token")
		t, err := app.tokenService.FindByID(tokenID)
		if err != nil {
			// no token or DB error
			// TODO: redirect user to proper page
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// validate token
		now := time.Now().UTC()
		if t.ExpiresAt.Before(now) {
			// token has expired
			// TODO: redirect user to proper page
			app.clientError(w, http.StatusBadRequest)
			return
		}
		if t.Kind != token.ValidateEmail {
			// improper token kind
			// TODO: redirect user to proper page
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// create session
		sid, err := app.sessionService.Create(t.UID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.AddCookie(w, *sid)

		// get user data
		u, err := app.userService.FindByID(t.UID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// activate account. Set u.Active = true
		err = app.userService.UpdateStatus(u.ID, true)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// update u.email_was_validated
		err = app.userService.UpdateEmailValidated(u.ID, true)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// delete current token
		err = app.tokenService.Delete(t.ID)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// At this point the user can not use this token again but,
		// as her/his account is already activated, in case she/he misses the
		// password creation it is possible to request a reset password

		// redirect user to create-password page
		tplData := &templateData{
			Title: "Create Password",
			User: &user.User{
				Name:  u.Name,
				Email: u.Email,
				ID:    u.ID,
			},
		}
		out, err := app.render("create-password.gohtml", tplData)
		if err != nil {
			app.serverError(w, err)
			return
		}

		out.WriteTo(w)
		if err != nil {
			app.serverError(w, err)
			return
		}
	})
}

// UpdateUser (PUT)
// func (app App) updateUser() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			app.clientError(w, http.StatusBadRequest)
// 			return
// 		}

// 		var data userData

// 		err := decodeJSONBody(w, r, &data)
// 		if err != nil {
// 			var mr *malformedRequest
// 			if errors.As(err, &mr) {
// 				http.Error(w, mr.msg, mr.status)
// 			} else {
// 				app.serverError(w, err)
// 			}
// 			return
// 		}
// 		if validationErrors := data.validate(); len(validationErrors) > 0 {
// 			err := map[string]interface{}{"validationError": validationErrors}
// 			w.Header().Set("Content-type", "application/json")
// 			app.clientError(w, http.StatusBadRequest)
// 			js, e := json.Marshal(err)
// 			if e != nil {
// 				app.serverError(w, e)
// 				return
// 			}
// 			w.Write(js)
// 			return
// 		}
// 	})
// 	// w.Write([]byte("updates a user's record"))
// }

// UpdateRoles (PUT)
// func (app App) updateRoles() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("updates a user's roles"))
// 	})
// }

// DeleteUser (DELETE)
// func (app App) DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	// w.Write([]byte("removes/deactivates a user account"))
// }
// GetUserByName (GET)
func (app App) getUserByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	// need a three-letter name at least
	if len(name) < 3 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	users, err := app.userService.FindByName(name)
	if err != nil {
		if errors.As(err, &user.ErrNoRecord) {
			w.Header().Set("Content-type", "application/json")
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	output, err := json.Marshal(users)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)

}

// GetUserByEmail (GET)
func (app App) getUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if !validEmail(email) {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	u, err := app.userService.FindByEmail(email)
	if err != nil {
		if errors.As(err, &user.ErrNoRecord) {
			w.Header().Set("Content-type", "application/json")
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	u.HashedPw = []byte("")
	output, err := json.Marshal(u)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (app App) getAllUsers(w http.ResponseWriter, r *http.Request) {
	var previous, next, pgSize string
	previous = r.URL.Query().Get("previous")
	next = r.URL.Query().Get("next")
	pgSize = r.URL.Query().Get("pgSize")
	size, err := strconv.Atoi(pgSize)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := app.userService.GetAll(previous, next, size)
	if err != nil {
		if err == user.ErrInvalidInputSyntax {
			app.clientError(w, http.StatusBadRequest)
		} else {
			app.serverError(w, err)
		}
		return
	}
	output, err := json.Marshal(res)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// resetPassword sends an reset password confirmation email with a "magic link"
// from: https://www.troyhunt.com/everything-you-ever-wanted-to-know/
func (app App) resetPassword() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// validate method
			if r.Method != http.MethodPost {
				app.clientError(w, http.StatusBadRequest)
				return
			}

			// validate email address
			var data struct {
				Email string `json:"email"`
			}
			err := decodeJSONBody(w, r, &data)
			if err != nil {
				var mr *malformedRequest
				if errors.As(err, &mr) {
					http.Error(w, mr.msg, mr.status)
				} else {
					app.serverError(w, err)
				}
				return
			}
			if !validEmail(data.Email) {
				app.clientError(w, http.StatusBadRequest)
				return
			}

			// does this account exists?
			u, err := app.userService.FindByEmail(data.Email)
			if err != nil {
				switch err {
				case user.ErrNoRecord:
					// TODO:
					// send msg to the given email address explaining that
					// someone tried to reset a password in this site, but
					// there is no such account...
					// redirect to default page explaining that an email
					// was sent to such account
					app.resetPasswordNoLocalUser(w, data.Email)
					app.redirectToSuccessPageForPwResetRequest(w, data.Email)
					return
				default:
					app.serverError(w, err)
					return
				}
			}

			// is this account active?
			if !u.Active || !u.EmailWasValidated {
				// TODO: redirect user to page explaining that he/she
				// needs to validate her/his email before changing the password
				app.clientError(w, http.StatusBadRequest)
				return
			}

			// delete all reset password tokens for this user
			allTks, err := app.tokenService.FindAllByUID(u.ID)
			if err == nil && allTks != nil {
				// delete tokens
				for _, t := range *allTks {
					if t.Kind == token.PwReset {
						err := app.tokenService.Delete(t.ID)
						if err != nil {
							app.serverError(w, err)
							return
						}
					}
				}
			} else if err != nil && err != token.ErrNoRecord {
				app.serverError(w, err)
				return
			}
			// create and store a new reset password token
			token, err := app.tokenService.New(u.ID, token.PwReset)
			if err != nil {
				app.serverError(w, err)
				return
			}
			// send msg/token to confirm user's email account  // compose email
			link := fmt.Sprintf("https://%s%s%s", r.Host, "/api/users/verifyResetPw?token=", *token)
			tplData := &templateData{
				Title: "Email Validation",
				User: &user.User{
					Name:  u.Name,
					Email: u.Email,
				},
				Link: link,
			}

			out, err := app.render("reset-password-request.gohtml", tplData)
			if err != nil {
				app.serverError(w, err)
				return
			}

			msg := out.String()
			s := "Reset de senha"
			res, err := app.mailerService.Send(u.Name, u.Email, "Reset de senha", msg)
			app.infoLog.Printf("Sent %v to %v\nResponse Code: %v\nMsg: %v\n",
				s, u.Email, res.Code, res.Msg)
			if err != nil {
				app.serverError(w, err)
				return
			}

			// redirect user to page explaining that an email was sent
			app.redirectToSuccessPageForPwResetRequest(w, u.Email)
			return
		})
}

func (app App) resetPasswordNoLocalUser(w http.ResponseWriter, email string) {
	// send msg to the given email address explaining that
	// someone tried to reset a password in this site, but
	// there is no such account...

	tplData := &templateData{
		Title: "Reset de senha",
		User: &user.User{
			Email: email,
		},
	}
	out, err := app.render("reset-pw-no-local-user-email.gohtml", tplData)
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := out.String()
	s := "Reset de Senha"
	res, err := app.mailerService.Send("", email, s, msg)
	app.infoLog.Printf("Sent %v to %v\nResponse Code: %v\nMsg: %v\n",
		s, email, res.Code, res.Msg)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app App) redirectToSuccessPageForPwResetRequest(w http.ResponseWriter, email string) {
	// redirect user to page explaining that an email was sent
	tplData := &templateData{
		Title: "Confirmação de Solicitação",
		User: &user.User{
			Email: email,
		},
	}
	out, err := app.render("default-reset-password-page.gohtml", tplData)
	if err != nil {
		app.serverError(w, err)
		return
	}
	_, err = out.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

//func (app App) getAll() http.Handler {
//return http.HandlerFunc(
//func(w http.ResponseWriter, r *http.Request) {
//Read body
//b, err := ioutil.ReadAll(r.Body)
//defer r.Body.Close()
//if err != nil {
//app.clientError(w, http.StatusBadRequest)
//return
//}

//var usersReq struct {
//Before string `json:"before"`
//After  string `json:"after"`
//PgSize int    `json:"pgsize"`
//}
//err = json.Unmarshal(b, &usersReq)
//if err != nil {
//app.serverError(w, err)
//return
//}

//res, err := app.userService.GetAll(usersReq.Before, usersReq.After, usersReq.PgSize)
//if err != nil {
//app.serverError(w, err)
//return
//}
//return values
//output, err := json.Marshal(res)
//if err != nil {
//app.serverError(w, err)
//return
//}
//w.Write(output)
//})
//}
func (app App) verifyResetPw() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			// extract token from query
			tokenID := r.URL.Query().Get("token")
			t, err := app.tokenService.FindByID(tokenID)
			if err != nil {
				// no token or DB error
				// TODO: redirect user to proper page
				msg := make(map[string]string)
				msg["Email validation"] = "No token"
				app.serverError(w, err)
				app.sendMsg(w, &msg)
				return
			}
			// validate token
			now := time.Now().UTC()
			if t.ExpiresAt.Before(now) {
				// token has expired
				// TODO: redirect user to proper page
				msg := make(map[string]string)
				msg["Email validation"] = "Expired token"
				app.serverError(w, token.ErrInvalidInputSyntax)
				app.sendMsg(w, &msg)
				return
			}
			if t.Kind != token.PwReset {
				// improper token
				// TODO: redirect user to proper page
				msg := make(map[string]string)
				msg["Email validation"] = "Invalid token"
				app.serverError(w, token.ErrInvalidInputSyntax)
				app.sendMsg(w, &msg)
				return
			}
			// create session
			sid, err := app.sessionService.Create(t.UID)
			if err != nil {
				app.serverError(w, err)
				return
			}
			app.AddCookie(w, *sid)

			// delete current token
			err = app.tokenService.Delete(t.ID)
			if err != nil {
				app.serverError(w, err)
				return
			}

			// get user data
			u, err := app.userService.FindByID(t.UID)
			if err != nil {
				app.serverError(w, err)
				return
			}

			// redirect user to change-password page
			tplData := &templateData{
				Title: "Create Password",
				User: &user.User{
					Name:  u.Name,
					Email: u.Email,
					ID:    u.ID,
				},
			}
			out, err := app.render("change-password.gohtml", tplData)
			if err != nil {
				app.serverError(w, err)
				return
			}
			_, err = out.WriteTo(w)
			if err != nil {
				app.serverError(w, err)
				return
			}
		})
}
