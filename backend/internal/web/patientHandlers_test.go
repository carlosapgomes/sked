package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
		{"Valid Search But Missing Patient", "Missing", 0, nil, http.StatusOK},
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
func TestGetPatientName(t *testing.T) {
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
		name     string
		searchID string
		wantBody []byte
		wantCode int
	}{
		{"Valid Search", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", []byte("Valid Patient"), http.StatusOK},
		{"Valid Search But Missing Patient", "2e134760-2006-4dc7-a315-025dc1081fb0", nil, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/patients/" + tt.searchID + "/name"
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
			var response patient.Patient
			if tt.wantCode == http.StatusOK {
				defer rs.Body.Close()
				respBody, _ := ioutil.ReadAll(rs.Body)
				//t.Logf("%s\n", respBody)
				err = json.Unmarshal(respBody, &response)
				if err != nil {
					t.Error("bad response body")
				}
			}
			if tt.wantBody != nil {
				if !bytes.Contains([]byte(response.Name), tt.wantBody) {
					t.Errorf("want body %s to contain %q", response, tt.wantBody)
				}
			}
		})
	}

}

func TestGetPatientPhones(t *testing.T) {
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
		name     string
		searchID string
		wantBody []byte
		wantCode int
	}{
		{"Valid Search", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", []byte("6544332135"), http.StatusOK},
		{"Valid Search But Missing Patient", "2e134760-2006-4dc7-a315-025dc1081fb0", nil, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/patients/" + tt.searchID + "/phones"
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
			var response struct {
				Phones []string
			}
			if tt.wantCode == http.StatusOK {
				defer rs.Body.Close()
				respBody, _ := ioutil.ReadAll(rs.Body)
				//t.Logf("%s\n", respBody)
				err = json.Unmarshal(respBody, &response)
				if err != nil {
					t.Error("bad response body")
				}
			}
			if tt.wantBody != nil {
				contains := false
				for _, e := range response.Phones {
					if bytes.Contains([]byte(e), tt.wantBody) {
						contains = true
					}
				}
				if !contains {
					t.Errorf("want body %s to contain %q", response, tt.wantBody)
				}

			}
		})
	}

}

func TestUpdatePatient(t *testing.T) {
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
		name      string
		patientID string
		newName   string
		newAdress string
		newCity   string
		newState  string
		newPhones []string
		wantCode  int
	}{
		{"Valid Update",
			"85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			"New Valid Patient Name",
			"New Address",
			"New City",
			"New State",
			[]string{"1234"},
			http.StatusOK,
		},
		{"Invalid Patient ID",
			"2e134760-2006-4dc7-a315-025dc1081fb0",
			"New Name",
			"New Address",
			"New City",
			"New State",
			[]string{"1234"},
			http.StatusInternalServerError,
		},
		{"Invalid Patient Name",
			"85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			"",
			"New Address",
			"New City",
			"New State",
			[]string{"1234"},
			http.StatusBadRequest,
		},
	}
	type putBody struct {
		ID      string   `json:"ID"`
		Name    string   `json:"Name"`
		Address string   `json:"Address"`
		City    string   `json:"City"`
		State   string   `json:"State"`
		Phones  []string `json:"Phones"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/patients/" + tt.patientID
			reqBody := &putBody{
				ID:      tt.patientID,
				Name:    tt.newName,
				Address: tt.newAdress,
				City:    tt.newCity,
				State:   tt.newState,
				Phones:  tt.newPhones,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			req, _ := http.NewRequest(http.MethodPut,
				path,
				strings.NewReader(string(body)))
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
			if tt.wantCode == http.StatusOK {
				path = ts.URL + "/patients?id=" + tt.patientID
				req, _ = http.NewRequest(http.MethodGet, path, nil)
				req.AddCookie(cookie)
				rs, err = ts.Client().Do(req)
				if err != nil {
					t.Fatal(err)
				}
				var response patient.Patient
				if rs.StatusCode == http.StatusOK {
					defer rs.Body.Close()
					respBody, _ := ioutil.ReadAll(rs.Body)
					err = json.Unmarshal(respBody, &response)
					if err != nil {
						t.Error("bad response body")
					}
					if (response.Name != tt.newName) ||
						(response.Address != tt.newAdress) ||
						(response.City != tt.newCity) ||
						(response.State != tt.newState) ||
						(response.Phones[0] != tt.newPhones[0]) {
						t.Error("Could not update patient")
					}
				}
			}

		})
	}

}

func TestUpdatePatientName(t *testing.T) {
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
		name      string
		patientID string
		newName   string
		wantCode  int
	}{
		{"Valid Update", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", "New Valid Patient Name", http.StatusOK},
		{"Invalid Patient ID", "2e134760-2006-4dc7-a315-025dc1081fb0", "New Name", http.StatusInternalServerError},
		{"Invalid Patient Name", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", "", http.StatusInternalServerError},
	}
	type postBody struct {
		Name string `json:"name"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/patients/" + tt.patientID + "/name"
			reqBody := &postBody{Name: tt.newName}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			req, _ := http.NewRequest(http.MethodPut, path, strings.NewReader(string(body)))
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
			if tt.wantCode == http.StatusOK {
				path = ts.URL + "/patients/" + tt.patientID + "/name"
				req, _ = http.NewRequest(http.MethodGet, path, nil)
				req.AddCookie(cookie)
				rs, err = ts.Client().Do(req)
				if err != nil {
					t.Fatal(err)
				}
				var response patient.Patient
				if rs.StatusCode == http.StatusOK {
					defer rs.Body.Close()
					respBody, _ := ioutil.ReadAll(rs.Body)
					//t.Logf("%s\n", respBody)
					err = json.Unmarshal(respBody, &response)
					if err != nil {
						t.Error("bad response body")
					}
					if response.Name != tt.newName {
						t.Errorf("Want new name %s but got %s", tt.newName, response.Name)
					}
				}
			}

		})
	}

}
func TestUpdatePatientPhones(t *testing.T) {
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
		name      string
		patientID string
		newPhones []string
		wantCode  int
	}{
		{"Valid Update", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", []string{"4321567"}, http.StatusOK},
		{"Invalid Patient ID", "2e134760-2006-4dc7-a315-025dc1081fb0", nil, http.StatusInternalServerError},
	}
	type postBody struct {
		Phones []string `json:"phones"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/patients/" + tt.patientID + "/phones"
			reqBody := &postBody{Phones: tt.newPhones}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			req, _ := http.NewRequest(http.MethodPut, path, strings.NewReader(string(body)))
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
			if tt.wantCode == http.StatusOK {
				path = ts.URL + "/patients/" + tt.patientID + "/phones"
				req, _ = http.NewRequest(http.MethodGet, path, nil)
				req.AddCookie(cookie)
				rs, err = ts.Client().Do(req)
				if err != nil {
					t.Fatal(err)
				}
				var response struct {
					Phones []string
				}
				if rs.StatusCode == http.StatusOK {
					defer rs.Body.Close()
					respBody, _ := ioutil.ReadAll(rs.Body)
					t.Logf("%s\n", respBody)
					err = json.Unmarshal(respBody, &response)
					if err != nil {
						t.Error("bad response body")
					}
					contains := false
					for _, e := range response.Phones {
						if bytes.Contains([]byte(e), []byte(tt.newPhones[0])) {
							contains = true
						}
					}
					if !contains {
						t.Errorf("want body %s to contain %q", response, tt.newPhones)
					}
				}
			}

		})
	}

}
func TestFindPatientByID(t *testing.T) {
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
		name     string
		searchID string
		wantBody []byte
		wantCode int
	}{
		{"Valid Search", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", []byte("Valid Patient"), http.StatusOK},
		{"Valid Search But Missing Patient", "2e134760-2006-4dc7-a315-025dc1081fb0", nil, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/patients?id=" + tt.searchID
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
			var response patient.Patient
			if tt.wantCode == http.StatusOK {
				defer rs.Body.Close()
				respBody, _ := ioutil.ReadAll(rs.Body)
				t.Logf("%s\n", respBody)
				err = json.Unmarshal(respBody, &response)
				if err != nil {
					t.Error("bad response body")
				}
			}
			if tt.wantBody != nil {
				if !bytes.Contains([]byte(response.Name), tt.wantBody) {
					t.Errorf("want body %s to contain %q", response, tt.wantBody)
				}
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

func TestGetAllPatients(t *testing.T) {
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
	testCases := []struct {
		desc          string
		previous      string
		next          string
		pgSize        int
		wantSize      int
		hasMore       bool
		wantCode      int
		wantContainID string
	}{
		{
			desc:          "Valid Page",
			previous:      "",
			next:          "",
			pgSize:        6,
			wantSize:      6,
			hasMore:       false,
			wantCode:      http.StatusOK,
			wantContainID: "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
		},
		{
			desc:          "Valid Cursor next",
			previous:      "NjhiMWQ1ZTItMzlkZC00NzEzLTg2MzEtYTA4MTAwMzgzYTBm",
			next:          "",
			pgSize:        2,
			wantSize:      1,
			hasMore:       false,
			wantCode:      http.StatusOK,
			wantContainID: "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
		},
		{
			desc:          "Valid Cursor previous",
			previous:      "",
			next:          "NjhiMWQ1ZTItMzlkZC00NzEzLTg2MzEtYTA4MTAwMzgzYTBm",
			pgSize:        2,
			wantSize:      2,
			hasMore:       true,
			wantCode:      http.StatusOK,
			wantContainID: "dcce1beb-aee6-4a4d-b724-94d470817323",
		},
	}
	// Page encapsulates data and pagination cursors
	type page struct {
		StartCursor     string            `json:"startCursor"`
		HasPreviousPage bool              `json:"hasPreviousPage"`
		EndCursor       string            `json:"endCursor"`
		HasNextPage     bool              `json:"hasNextPage"`
		Patients        []patient.Patient `json:"patients"`
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			path := ts.URL + "/patients?previous=" + tC.previous + "&next=" + tC.next + "&pgSize=" + strconv.Itoa(tC.pgSize)
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
			var cursor page
			if tC.wantCode == http.StatusOK {
				defer rs.Body.Close()
				respBody, _ := ioutil.ReadAll(rs.Body)
				//t.Logf("respBody: %s\n", respBody)
				err = json.Unmarshal(respBody, &cursor)
				if err != nil {
					t.Error("bad response body")
				}
			}
			//t.Logf("cursor: %v\n", cursor)
			code := rs.StatusCode
			if code != tC.wantCode {
				t.Errorf("Want %v; got %v\n", tC.wantCode, err)
			}
			if len(cursor.Patients) != tC.wantSize {
				t.Errorf("Want %v; got %v\n", tC.wantSize, len(cursor.Patients))
			}
			if tC.hasMore && !(cursor.HasNextPage || cursor.HasPreviousPage) {
				t.Errorf("want %v; got %v\n", tC.hasMore, (cursor.HasNextPage || cursor.HasPreviousPage))
			}
			var contain bool
			for _, p := range cursor.Patients {
				if p.ID == tC.wantContainID {
					contain = true
				}
			}
			if !contain {
				t.Errorf("Want response to contain %v ID;  but it did not\n", tC.wantContainID)
			}
		})
	}
}
