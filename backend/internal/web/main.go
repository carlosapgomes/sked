package web

import (
	"log"

	"carlosapgomes.com/gobackend/internal/mailer"
	"carlosapgomes.com/gobackend/internal/session"
	"carlosapgomes.com/gobackend/internal/token"
	"carlosapgomes.com/gobackend/internal/user"
)

type contextKey string

// ContextKeyIsAuthenticated key to set ctx as authenticated
const ContextKeyIsAuthenticated = contextKey("isAuthenticated")

// ContextKeySID key to set session ID in context
const ContextKeySID = contextKey("SID")

// ContextKeyUser key to set user data in context
const ContextKeyUser = contextKey("User")

// ContextResourceID holds the ID of the resources been accessed
const ContextResourceID = contextKey("resourceID")

// CkProps holds configurable cookie properties
type CkProps struct {
	Name     string
	HTTPOnly bool
	Secure   bool
	SameSite string
}

// App struct to hold its depedencies
type App struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	ckProps        *CkProps
	sessionService session.Service
	userService    user.Service
	mailerService  mailer.Service
	tokenService   token.Service
}

// New returns a new App object
func New(errorLog *log.Logger,
	infoLog *log.Logger,
	ckProps *CkProps,
	sessionService session.Service,
	userService user.Service,
	mailerService mailer.Service,
	tokenService token.Service) *App {
	return &App{
		errorLog:       errorLog,
		infoLog:        infoLog,
		ckProps:        ckProps,
		sessionService: sessionService,
		userService:    userService,
		mailerService:  mailerService,
		tokenService:   tokenService,
	}
}
