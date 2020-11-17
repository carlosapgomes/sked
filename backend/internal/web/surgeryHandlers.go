package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"carlosapgomes.com/sked/internal/surgery"
	"carlosapgomes.com/sked/internal/user"
)

type surgeriesData struct {
	ID              string `json:"id,omitempty"`
	DateTime        string `json:"dateTime"` // iso8601 format
	PatientName     string `json:"patientName"`
	PatientID       string `json:"patientID"`
	DoctorName      string `json:"doctorName"`
	DoctorID        string `json:"doctorID"`
	Notes           string `json:"notes"`
	ProposedSurgery string `json:"proposedSurgery"`
	CreatedBy       string `json:"createdBy"`
}

func (app App) surgeries() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/surgeries"):
			app.surgeriesNoPath(w, r)
		default:
			app.clientError(w, http.StatusBadRequest)
		}
	})
}

func (app App) surgeriesNoPath(w http.ResponseWriter, r *http.Request) {
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
			app.findSurgeryByID(w, r)
		case doctorID != "":
			app.findSurgeryByDoctorID(w, r)
		case patientID != "":
			app.findSurgeryByPatientID(w, r)
		case dt != "":
			app.findSurgeryByDate(w, r)
		case ((month != "") && (year != "")):
			app.findSurgeriesByMonth(w, r)
		default:
			app.getAllSurgeries(w, r)
		}
	case http.MethodPost:
		app.createSurgery(w, r)
	default:
		app.clientError(w, http.StatusBadRequest)
	}
}

func (app App) findSurgeryByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	surg, err := app.surgeryService.FindByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(*surg)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app App) findSurgeryByDoctorID(w http.ResponseWriter, r *http.Request) {
	doctorID := r.URL.Query().Get("doctorID")
	if doctorID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	surg, err := app.surgeryService.FindByDoctorID(doctorID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(surg)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app App) findSurgeryByPatientID(w http.ResponseWriter, r *http.Request) {
	patientID := r.URL.Query().Get("patientID")
	if patientID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	surg, err := app.surgeryService.FindByPatientID(patientID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(surg)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app App) findSurgeryByDate(w http.ResponseWriter, r *http.Request) {
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
	surg, err := app.surgeryService.FindByDate(dateTime)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(surg)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app App) findSurgeriesByMonth(w http.ResponseWriter, r *http.Request) {
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
	surgs, err := app.
		surgeryService.FindByMonthYear(m, y)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(surgs)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
func (app App) getAllSurgeries(w http.ResponseWriter, r *http.Request) {
	previous := r.URL.Query().Get("previous")
	next := r.URL.Query().Get("next")
	pgSize := r.URL.Query().Get("pgSize")
	size, err := strconv.Atoi(pgSize)
	if err != nil {
		app.serverError(w, err)
		return
	}
	//fmt.Printf("previous: %v\nnext: %v\npgSize: %v\n", previous, next, size)
	res, err := app.surgeryService.GetAll(previous, next, size)
	if err != nil {
		if err == surgery.ErrInvalidInputSyntax {
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

func (app App) createSurgery(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var newSurgery surgeriesData
	err = json.Unmarshal(b, &newSurgery)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if newSurgery.DateTime == "" ||
		newSurgery.PatientName == "" ||
		newSurgery.PatientID == "" ||
		newSurgery.DoctorName == "" ||
		newSurgery.DoctorID == "" ||
		newSurgery.Notes == "" ||
		newSurgery.ProposedSurgery == "" ||
		newSurgery.CreatedBy == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//dateTime, err := time.Parse("2006-01-02T15:04:05Z0700", newSurgery.DateTime)
	dateTime, err := time.Parse("2006-01-02T15:04:05-07:00", newSurgery.DateTime)
	if err != nil {
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
	id, err := app.surgeryService.Create(dateTime, newSurgery.PatientName,
		newSurgery.PatientID, newSurgery.DoctorName, newSurgery.DoctorID,
		newSurgery.Notes, newSurgery.ProposedSurgery, u.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	newSurgery.ID = *id
	output, err := json.Marshal(newSurgery)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
