package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
