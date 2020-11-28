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
	"carlosapgomes.com/sked/internal/services"
	"carlosapgomes.com/sked/internal/surgery"
	"carlosapgomes.com/sked/internal/web"
)

func TestFindSurgeryByID(t *testing.T) {
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
		nil,
		services.NewSurgeryService(mocks.NewSurgeryRepo(), userSvc),
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
		{"Valid Search But Missing Surgery", "2e134760-2006-4dc7-a315-025dc1081fb0", nil, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/surgeries?id=" + tt.searchID
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
			var response surgery.Surgery
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

func TestFindSurgeryByDoctorID(t *testing.T) {
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
		nil,
		services.NewSurgeryService(mocks.NewSurgeryRepo(), userSvc),
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
			path := ts.URL + "/surgeries?doctorID=" + tt.searchID
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
			var response []surgery.Surgery
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

func TestFindSurgeryByPatientID(t *testing.T) {
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
		nil,
		services.NewSurgeryService(mocks.NewSurgeryRepo(), userSvc),
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
			path := ts.URL + "/surgeries?patientID=" + tt.searchID
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
			var response []surgery.Surgery
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
func TestFindSurgeriesByMonth(t *testing.T) {
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
		nil,
		services.NewSurgeryService(mocks.NewSurgeryRepo(), userSvc),
	)
	ts := newTestServer(t, handlers.Routes())
	var tests = []struct {
		name     string
		month    int
		year     int
		wantSize int
		wantCode int
	}{
		{
			"Valid Search",
			9,
			2020,
			6,
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/surgeries?month=" + strconv.Itoa(tt.month) +
				"&year=" + strconv.Itoa(tt.year)
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
			var response []surgery.Surgery
			if tt.wantCode == http.StatusOK {
				defer rs.Body.Close()
				respBody, _ := ioutil.ReadAll(rs.Body)
				err = json.Unmarshal(respBody, &response)
				if err != nil {
					t.Error("bad response body")
				}
			}
			if len(response) != tt.wantSize {
				t.Errorf("Want response size %v, but got %v\n", tt.wantSize,
					len(response))
			}
		})
	}
}
func TestFindSurgeryByDate(t *testing.T) {
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
		nil,
		services.NewSurgeryService(mocks.NewSurgeryRepo(), userSvc),
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
		{"Valid Search But Missing Surgery on The Date", "2020-09-25", nil, http.StatusInternalServerError},
		{"Invalid Search", "202-09-25", nil, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := ts.URL + "/surgeries?date=" + tt.searchDate
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
			var response []surgery.Surgery
			if tt.wantCode == http.StatusOK {
				defer rs.Body.Close()
				respBody, _ := ioutil.ReadAll(rs.Body)
				t.Logf("%s\n", respBody)
				err = json.Unmarshal(respBody, &response)
				if err != nil {
					t.Error("bad response body")
				}
			}
			if tt.wantBody != nil && len(response) > 0 {
				if !bytes.Contains([]byte(response[0].PatientName), tt.wantBody) {
					t.Errorf("want body %s to contain %q", response[0].PatientName, tt.wantBody)
				}
			}
		})
	}
}

func TestCreateSurgery(t *testing.T) {

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
		nil,
		services.NewSurgeryService(mocks.NewSurgeryRepo(), userSvc),
	)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	var tests = []struct {
		name            string
		dateTime        string
		patientName     string
		patientID       string
		doctorName      string
		doctorID        string
		notes           string
		proposedSurgery string
		wantBody        []byte
		wantCode        int
	}{
		{"Valid Surgery", "2020-04-02T08:02:17-05:00",
			"John Doe", "c753a381-7642-4709-876f-57b16a5c6a6c", "Dr House",
			"f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes",
			"saphenectomy",
			[]byte("John Doe"), http.StatusOK},
	}
	type surgeriesData struct {
		DateTime        string `json:"dateTime"` // iso8601 format
		PatientName     string `json:"patientName"`
		PatientID       string `json:"patientID"`
		DoctorName      string `json:"doctorName"`
		DoctorID        string `json:"doctorID"`
		Notes           string `json:"notes"`
		ProposedSurgery string `json:"proposedSurgery"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := &surgeriesData{
				DateTime:        tt.dateTime,
				PatientName:     tt.patientName,
				PatientID:       tt.patientID,
				DoctorName:      tt.doctorName,
				DoctorID:        tt.doctorID,
				Notes:           tt.notes,
				ProposedSurgery: tt.proposedSurgery,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			t.Logf("%s\n", reqBody)
			path := ts.URL + "/surgeries"
			req, _ := http.NewRequest(http.MethodPost, path,
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
func TestGetAllSurgeries(t *testing.T) {
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
		nil,
		services.NewSurgeryService(mocks.NewSurgeryRepo(), userSvc),
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
			wantContainID: "7fef3c47-a01a-42a6-ac45-27a440596751",
		},
		{
			desc:          "Valid Cursor previous",
			previous:      "NWU2ZjdjZDEtZDhkMi00MGNkLTk3YTMtYWNhMDFhOTNiZmRl",
			next:          "",
			pgSize:        2,
			wantSize:      1,
			hasMore:       false,
			wantCode:      http.StatusOK,
			wantContainID: "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		},
		{
			desc:          "Valid Previous Cursor Bigger Response",
			previous:      "NzIzZTJmYTAtNzBhOS00YzIwLTg5ZDktYjVmNjk0MDViNzcy",
			next:          "",
			pgSize:        3,
			wantSize:      3,
			hasMore:       true,
			wantCode:      http.StatusOK,
			wantContainID: "7fef3c47-a01a-42a6-ac45-27a440596751",
		},
		{
			desc:          "Valid Next Cursor",
			previous:      "",
			next:          "NWU2ZjdjZDEtZDhkMi00MGNkLTk3YTMtYWNhMDFhOTNiZmRl",
			pgSize:        3,
			wantSize:      3,
			hasMore:       true,
			wantCode:      http.StatusOK,
			wantContainID: "19f66dc6-b5c8-497b-bba2-b982bb85ded8",
		},
	}
	// Page encapsulates data and pagination cursors
	type page struct {
		StartCursor     string            `json:"startCursor"`
		HasPreviousPage bool              `json:"hasPreviousPage"`
		EndCursor       string            `json:"endCursor"`
		HasNextPage     bool              `json:"hasNextPage"`
		Surgeries       []surgery.Surgery `json:"surgeries"`
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			path := ts.URL + "/surgeries?previous=" + tC.previous + "&next=" + tC.next + "&pgSize=" + strconv.Itoa(tC.pgSize)
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
			if len(cursor.Surgeries) != tC.wantSize {
				t.Errorf("Want %v; got %v\n", tC.wantSize, len(cursor.Surgeries))
			}
			if tC.hasMore && !(cursor.HasNextPage || cursor.HasPreviousPage) {
				t.Errorf("want %v; got %v\n", tC.hasMore, (cursor.HasNextPage || cursor.HasPreviousPage))
			}
			var contain bool
			for _, p := range cursor.Surgeries {
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
