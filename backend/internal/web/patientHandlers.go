package web

import (
	"net/http"
)

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

// GetAllPatients
func (app App) getAllPatients() http.Handler {
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
