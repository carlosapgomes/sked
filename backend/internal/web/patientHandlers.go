package web

import (
	"net/http"
	"strings"
)

func (app App) patients() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/patients"):
			app.patientsNoPath(w, r)
		case strings.HasSuffix(path, "/name"):
			app.patientName(w, r)
		case strings.HasSuffix(path, "/phones"):
			app.patientPhones(w, r)
		}
	})
}
func (app App) patientsNoPath(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//
	}
}
func (app App) patientName(w http.ResponseWriter, r *http.Request) {

}
func (app App) patientPhones(w http.ResponseWriter, r *http.Request) {

}

// CreatePatient creates a new patient record
func (app App) createPatient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

// FindPatientByID
func (app App) findPatientByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

// FindPatientByName
func (app App) findPatientByName() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

// GetAllPatients
func (app App) getAllPatients() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

// UpdatePatientName
func (app App) updatePatientName() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

// UpdatePatientPhone
func (app App) updatePatientPhone() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
