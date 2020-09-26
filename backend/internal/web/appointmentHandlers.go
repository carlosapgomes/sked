package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func (app App) appointments() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/appointments"):
			app.appointmentsNoPath(w, r)
		default:
			app.clientError(w, http.StatusBadRequest)
		}
	})
}

func (app App) appointmentsNoPath(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		doctorID := r.URL.Query().Get("doctorID")
		patientID := r.URL.Query().Get("patientID")
		dt := r.URL.Query().Get("date")
		switch {
		case id != "":
			app.findAppointmentByID(w, r)
		case doctorID != "":
			app.findAppointmentByDoctorID(w, r)
		case patientID != "":
			app.findAppointmentByPatientID(w, r)
		case dt != "":
			app.findAppointmentByDate(w, r)
		default:
			app.getAllAppointments(w, r)
		}
	case http.MethodPost:
		app.createAppointment(w, r)
	default:
		app.clientError(w, http.StatusBadRequest)
	}
}

func (app App) findAppointmentByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	appointmt, err := app.appointmentService.FindByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(*appointmt)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write(res)
}

func (app App) findAppointmentByDoctorID(w http.ResponseWriter, r *http.Request) {
	doctorID := r.URL.Query().Get("doctorID")
	if doctorID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	appointmt, err := app.appointmentService.FindByDoctorID(doctorID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(appointmt)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write(res)
}

func (app App) findAppointmentByPatientID(w http.ResponseWriter, r *http.Request) {
	patientID := r.URL.Query().Get("patientID")
	if patientID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	appointmt, err := app.appointmentService.FindByPatientID(patientID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(appointmt)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write(res)
}

func (app App) findAppointmentByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	dateTime, err := time.Parse("YYYY-MM-DD", date)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	appointmt, err := app.appointmentService.FindByDate(dateTime)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(appointmt)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write(res)
}

func (app App) getAllAppointments(w http.ResponseWriter, r *http.Request) {
}

func (app App) createAppointment(w http.ResponseWriter, r *http.Request) {
}
