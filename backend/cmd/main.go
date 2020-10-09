package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"carlosapgomes.com/sked/internal/services"
	"carlosapgomes.com/sked/internal/storage"
	"carlosapgomes.com/sked/internal/web"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// configuration precedence:
	// CLI > ENV > Defaults

	// First, set all the defaults:
	addr := ":9000"
	pgstr := "" // "postgres://sked:skedp@localhost/sked?sslmode=disable"
	ckName := "sid"
	ckSameSite := "strict" // Strict||Lax
	ckSecure := true
	ckHTTPOnly := true
	ssLifeTime := 20
	fromName := "Go Backend Manager"
	// create this email and register it in sendgrid
	// control panel as an authorized email sender
	fromAddress := "sked.manager@gmail.com"
	sgKey := ""

	// Lookup all corresponding Env vars
	if ad, ok := os.LookupEnv("HTTP_ADDR"); ok {
		addr = ad
	}
	if conn, ok := os.LookupEnv("PG_CONN"); ok {
		pgstr = conn
	}
	if cn, ok := os.LookupEnv("COOKIE_NAME"); ok {
		ckName = cn
	}
	if ssStr, ok := os.LookupEnv("COOKIE_SAMESITE"); ok {
		ckSameSite = ssStr
	}
	if csStr, ok := os.LookupEnv("COOKIE_SECURE"); ok {
		if cs, err := strconv.ParseBool(csStr); err == nil {
			ckSecure = cs
		}
	}
	if httoStr, ok := os.LookupEnv("COOKIE_HTTPONLY"); ok {
		if htto, err := strconv.ParseBool(httoStr); err == nil {
			ckHTTPOnly = htto
		}
	}
	if slStr, ok := os.LookupEnv("SESSION_LIFETIME"); ok {
		if sl, err := strconv.ParseInt(slStr, 10, 0); err == nil {
			ssLifeTime = int(sl)
		}
	}
	if key, ok := os.LookupEnv("SENDGRID_API_KEY"); ok {
		sgKey = key
	}
	if name, ok := os.LookupEnv("FROM_NAME"); ok {
		fromName = name
	}
	if add, ok := os.LookupEnv("FROM_EMAIL"); ok {
		fromAddress = add
	}

	// Look for all CLI flags, overwriting what was set until now if a flag is provided
	flag.StringVar(&addr, "addr", addr, "HTTP network address")
	flag.StringVar(&pgstr, "pgstr", pgstr, "Postgres data source name")
	flag.StringVar(&ckName, "ckname", ckName, "cookie session name (default: 'sid')")
	flag.StringVar(&ckSameSite, "ckss", ckSameSite, "cookie same site (default: 'true')")
	flag.BoolVar(&ckSecure, "cksec", ckSecure, "cookie secure (default: 'true')")
	flag.BoolVar(&ckHTTPOnly, "ckonly", ckHTTPOnly, "cookie http only (default: 'true')")
	flag.IntVar(&ssLifeTime, "slife", ssLifeTime, "session life time in minutes(default: '20')")
	flag.StringVar(&sgKey, "sgkey", sgKey, "Sendgrid API key")
	flag.StringVar(&fromName, "from", fromName, "'From' name to use in email")
	flag.StringVar(&fromAddress, "email", fromAddress, "'From' email to use")
	flag.Parse()

	// Finallly, check for external or secret vars that are
	// required not to be empty
	if pgstr == "" {
		flag.PrintDefaults()
		errorLog.Fatal(errors.New("Please, supply DB connection string"))
		os.Exit(1)
	}

	if sgKey == "" {
		flag.PrintDefaults()
		errorLog.Fatal(errors.New("Please, supply a Sendgrid API key"))
		os.Exit(1)
	}

	db, err := storage.NewDB(pgstr)
	if err != nil {
		errorLog.Panic(err)
	}
	// defer a call to db.Close(), so that the connection pool is
	// closed before the main() function exits.
	defer db.Close()

	// Initialize repositories instances (driven adapters)
	userRepository := storage.NewPgUserRepository(db)
	sessionRepository := storage.NewPgSessionRepository(db)
	tokenRepository := storage.NewPgTokenRepository(db)
	patientRepository := storage.NewPgPatientRepository(db)
	appointmentRepository := storage.NewPgAppointmentRepository(db)
	surgeryRepository := storage.NewPgSurgeryRepository(db)

	// Initialize core services injecting its dependencies
	// when neccessary
	userService := services.NewUserService(userRepository)
	sessionService := services.NewSessionService(ssLifeTime, sessionRepository)
	mailerService := services.NewMailerService(sgKey, fromName, fromAddress)
	tokenService := services.NewTokenService(tokenRepository)
	patientService := services.NewPatientService(patientRepository)
	appointmentService := services.NewAppointmentService(appointmentRepository,
		userService)
	surgeryService := services.NewSurgeryService(surgeryRepository,
		userService)
	ckprops := &web.CkProps{
		Name:     ckName,
		HTTPOnly: ckHTTPOnly,
		Secure:   ckSecure,
		SameSite: ckSameSite,
	}
	// Initialize web adapter (driver adapter) injecting the
	// core services as its dependencies
	app := web.New(errorLog,
		infoLog,
		ckprops,
		sessionService,
		userService,
		mailerService,
		tokenService,
		patientService,
		appointmentService,
		surgeryService)

	// Run web adapter
	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
		// Add Idle, Read and Write timeouts to the server.
		// should modify this and create a per handler timeout if
		// I want to upload/download files
		// see: https://medium.com/@simonfrey/go-as-in-golang-standard-net-http-config-will-break-your-production-environment-1360871cb72b
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
