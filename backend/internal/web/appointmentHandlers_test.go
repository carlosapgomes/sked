package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"carlosapgomes.com/sked/internal/appointment"
	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/services"
	"carlosapgomes.com/sked/internal/web"
)

func TestFindAppointmentByID(t *testing.T) {
	userSvc := services.NewUserService(mocks.NewUserRepo())
	handlers := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		userSvc,
		nil,
		mocks.NewTokenMockSvc(),
		nil,
		services.NewAppointmentService(mocks.NewAppointmentRepo(), userSvc),
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
		{"Valid Search", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", []byte("John Doe"), http.StatusOK},
		{"Valid Search But Missing Appointment", "2e134760-2006-4dc7-a315-025dc1081fb0", nil, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/appointments?id=" + tt.searchID
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
			var response appointment.Appointment
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
				if !bytes.Contains([]byte(response.PatientName), tt.wantBody) {
					t.Errorf("want body %s to contain %q", response.PatientName, tt.wantBody)
				}
			}
		})
	}
}

func TestFindAppointmentByDoctorID(t *testing.T) {
	userSvc := services.NewUserService(mocks.NewUserRepo())
	handlers := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		userSvc,
		nil,
		mocks.NewTokenMockSvc(),
		nil,
		services.NewAppointmentService(mocks.NewAppointmentRepo(), userSvc),
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
		{"Valid Search", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", []byte("John Doe"), http.StatusOK},
		{"Valid Search But Missing DoctorID", "2e134760-2006-4dc7-a315-025dc1081fb0", nil, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/appointments?doctorID=" + tt.searchID
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
			var response []appointment.Appointment
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
				if !bytes.Contains([]byte(response[0].PatientName), tt.wantBody) {
					t.Errorf("want body %s to contain %q", response[0].PatientName, tt.wantBody)
				}
			}
		})
	}
}

func TestFindAppointmentByPatientID(t *testing.T) {
	userSvc := services.NewUserService(mocks.NewUserRepo())
	handlers := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		userSvc,
		nil,
		mocks.NewTokenMockSvc(),
		nil,
		services.NewAppointmentService(mocks.NewAppointmentRepo(), userSvc),
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
		{"Valid Search", "22070f56-5d52-43f0-9f59-5de61c1db506", []byte("John Doe"), http.StatusOK},
		{"Valid Search But Missing PatientID", "2e134760-2006-4dc7-a315-025dc1081fb0", nil, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/appointments?patientID=" + tt.searchID
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
			var response []appointment.Appointment
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
				if !bytes.Contains([]byte(response[0].PatientName), tt.wantBody) {
					t.Errorf("want body %s to contain %q", response[0].PatientName, tt.wantBody)
				}
			}
		})
	}
}
func TestFindAppointmentByDate(t *testing.T) {
	userSvc := services.NewUserService(mocks.NewUserRepo())
	handlers := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		userSvc,
		nil,
		mocks.NewTokenMockSvc(),
		nil,
		services.NewAppointmentService(mocks.NewAppointmentRepo(), userSvc),
		nil,
	)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	var tests = []struct {
		name       string
		searchDate string
		wantBody   []byte
		wantCode   int
	}{
		{"Valid Search", "2020-09-06", []byte("John Doe"), http.StatusOK},
		{"Valid Search But Missing Appointment on The Date", "2020-09-25", nil, http.StatusInternalServerError},
		{"Invalid Search", "202-09-25", nil, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/appointments?date=" + tt.searchDate
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
			var response []appointment.Appointment
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
				if !bytes.Contains([]byte(response[0].PatientName), tt.wantBody) {
					t.Errorf("want body %s to contain %q", response[0].PatientName, tt.wantBody)
				}
			}
		})
	}
}

func TestCreateAppointment(t *testing.T) {

	userSvc := services.NewUserService(mocks.NewUserRepo())
	handlers := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		userSvc,
		nil,
		mocks.NewTokenMockSvc(),
		nil,
		services.NewAppointmentService(mocks.NewAppointmentRepo(), userSvc),
		nil,
	)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	var tests = []struct {
		name        string
		dateTime    string
		patientName string
		patientID   string
		doctorName  string
		doctorID    string
		notes       string
		createdBy   string
		wantBody    []byte
		wantCode    int
	}{
		{"Valid Appointment", "2020-04-02T08:02:17-05:00", "John Doe", "c753a381-7642-4709-876f-57b16a5c6a6c", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "896d45e7-b544-41da-aa3f-f59a321fcdb9", []byte("John Doe"), http.StatusOK},
	}
	type appointmentsData struct {
		DateTime    string `json:"dateTime"` // iso8601 format
		PatientName string `json:"patientName"`
		PatientID   string `json:"patientID"`
		DoctorName  string `json:"doctorName"`
		DoctorID    string `json:"doctorID"`
		Notes       string `json:"notes"`
		CreatedBy   string `json:"createdBy"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := &appointmentsData{
				DateTime:    tt.dateTime,
				PatientName: tt.patientName,
				PatientID:   tt.patientID,
				DoctorName:  tt.doctorName,
				DoctorID:    tt.doctorID,
				Notes:       tt.notes,
				CreatedBy:   tt.createdBy,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			t.Logf("%s\n", reqBody)
			path := ts.URL + "/appointments"
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
			if tt.wantCode == http.StatusOK && tt.wantBody != nil {
				defer rs.Body.Close()
				respBody, _ := ioutil.ReadAll(rs.Body)
				//t.Logf("%s\n", respBody)
				if !bytes.Contains([]byte(respBody), tt.wantBody) {
					t.Errorf("want body %s to contain %q\n", respBody, tt.wantBody)
				}
			}
		})
	}
}
