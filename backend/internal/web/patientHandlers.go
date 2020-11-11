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
		case strings.HasPrefix(path, "/patients/"):
			app.updatePatient(w, r)
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
			app.findPatientByID(w, r)
		case name != "":
			app.findPatientByName(w, r)
		default:
			app.getAllPatients(w, r)
		}
	case http.MethodPost:
		app.createPatient(w, r)
	default:
		app.clientError(w, http.StatusBadRequest)
	}
}
func (app App) updatePatient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		app.clientError(w, http.StatusBadRequest)
	}
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var p patientData
	err = json.Unmarshal(b, &p)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if validationErrors := p.validate(); len(validationErrors) > 0 {
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
	err = app.patientService.UpdatePatient(
		p.ID,
		p.Name,
		p.Address,
		p.City,
		p.State,
		p.Phones,
		u.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write()
}
func (app App) patientName(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
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
	res, err := json.Marshal(&map[string]string{"Name": p.Name})
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
		Name string `json:"name"`
	}
	err = json.Unmarshal(b, &newName)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// get logged in user from ctx
	u, ok := r.Context().Value(ContextKeyUser).(*user.User)
	if !ok {
		// no ContextKeyUser -> user is not authenticated
		app.clientError(w, http.StatusForbidden)
		return
	}
	err = app.patientService.UpdateName(pID, newName.Name, u.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app App) patientPhones(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getPatientPhones(w, r)
	case http.MethodPut:
		app.updatePatientPhones(w, r)
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

// getPatientPhones returns a list of a patient's phone numbers
func (app App) getPatientPhones(w http.ResponseWriter, r *http.Request) {
	pID := app.between(r.URL.Path, "/patients/", "/phones")
	if pID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	p, err := app.patientService.FindByID(pID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(&map[string][]string{"Phones": p.Phones})
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// UpdatePatientPhone, updates a patient's phones list
// it expects an json in the body like
// {"Phones":["1234","34534"]}
func (app App) updatePatientPhones(w http.ResponseWriter, r *http.Request) {
	pID := app.between(r.URL.Path, "/patients/", "/phones")
	if pID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var list struct {
		Phones []string `json:"Phones"`
	}
	err := decodeJSONBody(w, r, &list)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// get logged in user from ctx
	u, ok := r.Context().Value(ContextKeyUser).(*user.User)
	if !ok {
		// no ContextKeyUser -> user is not authenticated
		app.clientError(w, http.StatusForbidden)
		return
	}
	err = app.patientService.UpdatePhone(pID, list.Phones, u.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	id, err = app.patientService.Create(newPatient.Name, newPatient.Address,
		newPatient.City, newPatient.State, newPatient.Phones, u.ID)
	if err != nil {
		if errors.As(err, &patient.ErrDuplicateField) {
			w.Header().Set("Content-type", "application/json")
			app.clientError(w, http.StatusBadRequest)
			js, err := json.Marshal(&map[string]string{
				"Name": "there is a patient record with that name",
			})
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// FindPatientByID finds a patient with the given ID
// on the query's params
// patient?id=xxxxxxxx
// and return a json with the patient data
func (app App) findPatientByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	p, err := app.patientService.FindByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := json.Marshal(*p)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// getAllPatients
func (app App) getAllPatients(w http.ResponseWriter, r *http.Request) {
	previous := r.URL.Query().Get("previous")
	next := r.URL.Query().Get("next")
	pgSize := r.URL.Query().Get("pgSize")
	if pgSize == "" {
		// set a mininum page size
		pgSize = "5"
	}
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
