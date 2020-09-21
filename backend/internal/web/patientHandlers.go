package web

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"carlosapgomes.com/sked/internal/patient"
	"carlosapgomes.com/sked/internal/user"
)

type patientData struct {
	ID      string   `json:"ID,omitempty"`
	Name    string   `json:"Name"`
	Address string   `json:"Address"`
	City    string   `json:"City"`
	State   string   `json:"State"`
	Phones  []string `json:"Phones"`
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
		errs.Add("Name", "The name field must be between 3-100 chars!")
	}

	// check if phones list is empty
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
	case http.MethodPost:
		app.createPatient(w, r)
	default:
		app.clientError(w, http.StatusBadRequest)
	}
}
func (app App) patientName(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.getPatientName(w, r)
	case http.MethodPut:
		app.updatePatientName(w, r)
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}
func (app App) getPatientName(w http.ResponseWriter, r *http.Request) {
	pID := app.between(r.URL.Path, "/patients/", "/name")
	if pID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	p, err := app.patientService.FindByID(pID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(&map[string]string{"name": p.Name})
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write(res)

}

// updatePatientName expects a json in the request body
// like: {"name":"This Is The New Name"}
func (app App) updatePatientName(w http.ResponseWriter, r *http.Request) {

	pID := app.between(r.URL.Path, "/patients/", "/name")
	if pID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var newName struct {
		name string
	}
	err = json.Unmarshal(b, &newName)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.patientService.UpdateName(pID, newName)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	if validationErrors := newPatient.validate(); len(validationErrors) > 0 {
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
	// get logged in user from ctx
	u, ok := r.Context().Value(ContextKeyUser).(*user.User)
	if !ok {
		// no ContextKeyUser -> user is not authenticated
		app.clientError(w, http.StatusForbidden)
		return
	}
	var id *string
	id, err = app.patientService.Create(newPatient.Name, newPatient.Address, newPatient.City, newPatient.State, newPatient.Phones, u.ID)
	if err != nil {
		if errors.As(err, &patient.ErrDuplicateField) {
			w.Header().Set("Content-type", "application/json")
			app.clientError(w, http.StatusBadRequest)
			js, err := json.Marshal(&map[string]string{"Name": "there is a patient record with that name"})
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
	newPatient.ID = *id
	output, err := json.Marshal(newPatient)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write(output)
}

// FindPatientByID
func (app App) findPatientByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))

}

// FindPatientByName find patient(s) with the name given
// on the query's params
// /patients?name=name_to_find
// and returns a json with an array of patients
func (app App) findPatientByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	list, err := app.patientService.FindByName(name)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(*list)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write(res)
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

// UpdatePatientPhone
func (app App) updatePatientPhone() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
