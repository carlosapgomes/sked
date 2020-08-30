package web

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

	"carlosapgomes.com/sked/internal/user"
	"github.com/golang/gddo/httputil/header"
	uuid "github.com/satori/go.uuid"
)

// stdChars is a set of standard characters allowed in uniuri string.
var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app App) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	//fmt.Printf("Error from serverError: %s", err)
	//fmt.Printf("stack trace: \n%s\n", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app App) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app App) sendMsg(w http.ResponseWriter, msg *map[string]string) {
	js, err := json.Marshal(msg)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app App) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ContextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app App) isAdmin(r *http.Request) bool {
	u, ok := r.Context().Value(ContextKeyUser).(*user.User)
	if !ok {
		app.infoLog.Print("isAdmin: user data is not in ctx")
		return false
	}
	for _, role := range u.Roles {
		if role == user.RoleAdmin {
			app.infoLog.Print("isAdmin: requesting user is Admin")
			return true
		}
	}
	app.infoLog.Print("isAdmin: requesting user is NOT Admin")
	return false
}
func (app App) findRole(roles []string, r string) bool {
	for _, role := range roles {
		if role == r {
			return true
		}
	}
	return false
}
func (app App) isValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}

// find return substring between two other strings
func (app App) between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

// IsSameUserOrAdmin tests if current user is trying to access his/her own data or is admin
func (app App) IsSameUserOrAdmin(r *http.Request) bool {
	// try to find some context data
	ctxuID, ok := r.Context().Value(ContextResourceID).(string)
	if !ok {
		ctxuID = ""
	}
	u, ok := r.Context().Value(ContextKeyUser).(*user.User)
	if !ok {
		app.infoLog.Print("IsSameUserOrAdmin: user is not authenticated")
		return false
	}
	if r.Method == http.MethodGet {
		email := r.URL.Query().Get("email")
		uid := r.URL.Query().Get("uid")
		app.infoLog.Printf("IsSameUserOrAdmin u.ID: %s", u.ID)
		app.infoLog.Printf("IsSameUserOrAdmin ctxuID: %s", ctxuID)
		if (email == u.Email) || (uid == u.ID) || ((ctxuID != "") && (ctxuID == u.ID)) {
			app.infoLog.Print("IsSameUserOrAdmin: requesting user GET is the same")
			return true
		}
		app.infoLog.Print("IsSameUserOrAdmin: requesting user GET is NOT the same")
	} else {
		// for POST,PUT,PATCH,DEL methods, data must be in body
		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			var data map[string]interface{}
			err = json.Unmarshal(b, &data)
			if err == nil {
				if e, ok := data["email"]; ok {
					if e.(string) == u.Email {
						return true
					}
				} else {
					if id, ok := data["uid"]; ok {
						if id.(string) == u.ID {
							app.infoLog.Print("IsSameUserOrAdmin: requesting user non-GET is the same")
							return true
						}
						app.infoLog.Print("IsSameUserOrAdmin: requesting user non-GET is NOT the same")
					}
				}
			}
		}
	}
	return app.isAdmin(r)
}

// From:
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

// generateToken generate a 40 characters random token
// based on: https://github.com/dchest/uniuri
func generateToken() string {
	return newLenChars(40, stdChars)
}

func newLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("uniuri: wrong charset length for NewLenChars")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("uniuri: error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
