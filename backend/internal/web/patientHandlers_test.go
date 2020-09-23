package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/patient"
	"carlosapgomes.com/sked/internal/services"
	"carlosapgomes.com/sked/internal/web"
)

func TestFindPatientByName(t *testing.T) {
	handlers := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		services.NewUserService(mocks.NewUserRepo()),
		nil,
		mocks.NewTokenMockSvc(),
		services.NewPatientService(mocks.NewPatientRepo()),
		nil,
		nil,
	)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	var tests = []struct {
		name       string
		searchName string
		wantSize   int
		wantBody   []byte
		wantCode   int
	}{
		{"Valid Search", "Valid", 1, []byte("Valid Patient"), http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/patients?name=" + tt.searchName
			req, _ := http.NewRequest(http.MethodGet, path, nil)
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "167ced64-af16-45d2-bb08-e35233c04ad1",
			}
			req.AddCookie(cookie)
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			code := rs.StatusCode
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			var response []patient.Patient
			defer rs.Body.Close()
			respBody, _ := ioutil.ReadAll(rs.Body)
			err = json.Unmarshal(respBody, &response)
			if err != nil {
				t.Error("bad response body")
			}
			if len(response) != tt.wantSize {
				t.Errorf("want response size %d; got %d", tt.wantSize, len(response))
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}
}
func TestCreatePatient(t *testing.T) {
	handlers := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		services.NewUserService(mocks.NewUserRepo()),
		nil,
		mocks.NewTokenMockSvc(),
		services.NewPatientService(mocks.NewPatientRepo()),
		nil,
		nil,
	)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	var tests = []struct {
		name        string
		patientName string
		address     string
		city        string
		state       string
		phones      []string
		wantCode    int
	}{
		{"Valid New Patient", "Valid patient", "Main Street 34", "Capital City", "TH", []string{"123456"}, http.StatusOK},
		{"Empty name", "", "Main Street 34", "Capital City", "TH", []string{"123456"}, http.StatusBadRequest},
		{"Invalid name length", "no", "Main Street 34", "Capital City", "TH", []string{"123456"}, http.StatusBadRequest},
	}
	type patientData struct {
		ID      string   `json:"ID,omitempty"`
		Name    string   `json:"Name"`
		Address string   `json:"Address"`
		City    string   `json:"City"`
		State   string   `json:"State"`
		Phones  []string `json:"Phones"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := &patientData{
				Name:    tt.patientName,
				Address: tt.address,
				City:    tt.city,
				State:   tt.state,
				Phones:  tt.phones,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			path := ts.URL + "/patients"
			req, _ := http.NewRequest(http.MethodPost, path, strings.NewReader(string(body)))
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "167ced64-af16-45d2-bb08-e35233c04ad1",
			}
			req.AddCookie(cookie)
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			code := rs.StatusCode
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

		})
	}
}
