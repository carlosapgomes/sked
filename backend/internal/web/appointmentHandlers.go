package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"carlosapgomes.com/sked/internal/appointment"
	"carlosapgomes.com/sked/internal/user"
)

type appointmentsData struct {
	ID          string `json:"id,omitempty"`
	DateTime    string `json:"dateTime"` // iso8601 format
	PatientName string `json:"patientName"`
	PatientID   string `json:"patientID"`
	DoctorName  string `json:"doctorName"`
	DoctorID    string `json:"doctorID"`
	Notes       string `json:"notes"`
	CreatedBy   string `json:"createdBy"`
}

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
		month := r.URL.Query().Get("month")
		year := r.URL.Query().Get("year")
		switch {
		case id != "":
			app.findAppointmentByID(w, r)
		case doctorID != "":
			app.findAppointmentByDoctorID(w, r)
		case patientID != "":
			app.findAppointmentByPatientID(w, r)
		case dt != "":
			app.findAppointmentByDate(w, r)
		case ((month != "") && (year != "")):
			app.findAppointmentsByMonth(w, r)
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
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app App) findAppointmentByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	dateTime, err := time.Parse("2006-01-02", date)
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app App) findAppointmentsByMonth(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	if (month == "") || (year == "") {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	m, err := strconv.Atoi(month)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(year)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	appointmts, err := app.
		appointmentService.FindByMonthYear(m, y)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(appointmts)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app App) getAllAppointments(w http.ResponseWriter, r *http.Request) {
	previous := r.URL.Query().Get("previous")
	next := r.URL.Query().Get("next")
	pgSize := r.URL.Query().Get("pgSize")
	size, err := strconv.Atoi(pgSize)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := app.appointmentService.GetAll(previous, next, size)
	if err != nil {
		if err == appointment.ErrInvalidInputSyntax {
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

func (app App) createAppointment(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var newAppointmt appointmentsData
	err = json.Unmarshal(b, &newAppointmt)
	if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Println(newAppointmt)
	if newAppointmt.DateTime == "" ||
		newAppointmt.PatientName == "" ||
		newAppointmt.PatientID == "" ||
		newAppointmt.DoctorName == "" ||
		newAppointmt.DoctorID == "" ||
		newAppointmt.Notes == "" {
		fmt.Println("some empty field")
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//dateTime, err := time.Parse("2006-01-02T15:04:05Z0700", newAppointmt.DateTime)
	dateTime, err := time.Parse("2006-01-02T15:04:05-07:00", newAppointmt.DateTime)
	if err != nil {
		fmt.Println("bad date format")
		app.clientError(w, http.StatusBadRequest)
		return
	}
	fmt.Println(dateTime.Format("2006-01-02T15:04:05-07:00"))

	// get logged in user from ctx
	u, ok := r.Context().Value(ContextKeyUser).(*user.User)
	if !ok {
		// no ContextKeyUser -> user is not authenticated
		app.clientError(w, http.StatusForbidden)
		return
	}
	id, err := app.appointmentService.Create(
		dateTime,
		newAppointmt.PatientName,
		newAppointmt.PatientID,
		newAppointmt.DoctorName,
		newAppointmt.DoctorID,
		newAppointmt.Notes,
		u.ID)
	if err != nil {
		fmt.Printf("appointmentService error: %v\n", err)
		app.serverError(w, err)
		return
	}
	newAppointmt.ID = *id
	output, err := json.Marshal(newAppointmt)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
