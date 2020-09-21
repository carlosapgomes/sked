package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"carlosapgomes.com/sked/internal/patient"
)

type patientData struct {
	ID       string   `json:"ID,omitempty"`
	Name     string   `json:"Name"`
	Address  string   `json:"Address"`
	City     string   `json:"City"`
	State    string   `json:"State"`
	Phones   []string `json:"Phones"`
	Email    string   `json:"Email"`
	Phone    string   `json:"Phone"`
	Password string   `json:"Password"`
}

// validates request patient data
func (p *patientData) validate() url.Values {
	errs := url.Values{}

	// check if name empty
	if p.Name == "" {
		errs.Add("Name", "This field cannot be blank")
	}

	nameLen := utf8.RuneCountInString(p.Name)
	if nameLen < 3 || nameLen > 100 {
		errs.Add("name", "The name field must be between 3-100 chars!")
	}

	// check if email is empty
	if len(p.Phones) == 0 {
		errs.Add("Phones", "This field cannot be empty")
	}
	return errs
}

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
		default:
			app.clientError(w, http.StatusBadRequest)
		}
	})
}
func (app App) patientsNoPath(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")
		switch {
		case id != "":
			app.findPatientByName(w, r)
		case name != "":
			app.findPatientByID(w, r)
		default:
			app.getAllPatients(w, r)
		}
	case http.Post:
		app.createPatient(w, r)
	default:
		app.clientError(w, http.StatusBadRequest)
	}
}
func (app App) patientName(w http.ResponseWriter, r *http.Request) {

}
func (app App) patientPhones(w http.ResponseWriter, r *http.Request) {

}

// createPatient creates a new patient record
func (app App) createPatient(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var newPatient patientData
	err = json.Unmarshal(b, &newPatient)
	if err != nil {
		app.serverError(w, err)
		return
	}

}

// FindPatientByID
func (app App) findPatientByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))

}

// FindPatientByName
func (app App) findPatientByName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))

}

// getAllPatients
func (app App) getAllPatients(w http.ResponseWriter, r *http.Request) {
	previous := r.URL.Query().Get("previous")
	next := r.URL.Query().Get("next")
	pgSize := r.URL.Query().Get("pgSize")
	size, err := strconv.Atoi(pgSize)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := app.patientService.GetAll(previous, next, size)
	if err != nil {
		if err == patient.ErrInvalidInputSyntax {
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
	w.Write(output)
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
