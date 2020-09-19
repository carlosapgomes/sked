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
	"carlosapgomes.com/sked/internal/services"
	"carlosapgomes.com/sked/internal/user"
	"carlosapgomes.com/sked/internal/web"
)

// TestHealthz function implements health check (kind of unit test)
func TestHealthz(t *testing.T) {
	handlers := newTestApplication(t)
	// Notice that we defer a call to ts.Close() to shutdown the server when
	// the test finishes.
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	code, _, body := ts.get(t, "/healthz")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

// // TestUpdateUserStatus function tests user status update
// func TestUpdateUserStatus(t *testing.T) {
// 	app := newTestApplication(t)
// 	ts := newTestServer(t, app.Routes())
// 	defer ts.Close()
// 	tests := []struct {
// 		name     string
// 		ID       string
// 		Status   bool
// 		wantCode int
// 		wantBody []byte
// 	}{
// 		{"Valid activate account", "1", true, http.StatusOK, nil},
// 		{"Valid deactivate account", "1", false, http.StatusOK, nil},
// 	}
// 	type postBody struct {
// 		ID     string
// 		Status bool
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			reqBody := &postBody{ID: tt.ID,
// 				Status: tt.Status}
// 			body, err := json.Marshal(reqBody)
// 			if err != nil {
// 				t.Log(err)
// 			}
// 			code, _, respBody := ts.post(t, "/updateStatus", "applicaton/json", strings.NewReader(string(body)))
// 			if code != tt.wantCode {
// 				t.Errorf("want %d; got %d", tt.wantCode, code)
// 			}
// 			if !bytes.Contains(respBody, tt.wantBody) {
// 				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
// 			}
// 		})
// 	}
// 	type postBody2 struct {
// 		ID     string
// 		Status string
// 	}
// 	tests2 := []struct {
// 		name     string
// 		ID       string
// 		Status   string
// 		wantCode int
// 		wantBody []byte
// 	}{
// 		{"Valid string activate account", "1", "true", http.StatusOK, nil},
// 		{"Valid string deactivate account", "1", "false", http.StatusOK, nil},
// 		{"Invalid string deactivate account", "1", "", http.StatusBadRequest, nil},
// 	}
// 	for _, tt := range tests2 {
// 		t.Run(tt.name, func(t *testing.T) {
// 			reqBody := &postBody2{ID: tt.ID,
// 				Status: tt.Status}
// 			body, err := json.Marshal(reqBody)
// 			if err != nil {
// 				t.Log(err)
// 			}
// 			code, _, respBody := ts.post(t, "/updateStatus", "applicaton/json", strings.NewReader(string(body)))
// 			if code != tt.wantCode {
// 				t.Errorf("want %d; got %d", tt.wantCode, code)
// 			}
// 			if !bytes.Contains(respBody, tt.wantBody) {
// 				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
// 			}
// 		})
// 	}
// }

// TestValidateEmail function tests /users/validateEmail handler
func TestValidateEmail(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()
	tests := []struct {
		name     string
		token    string
		method   string
		wantCode int
		wantBody []byte
	}{
		{"Valid token", "Sqi1Z86I9632BndMkS632BNO4i1Z83v9KGd27FXs", http.MethodGet, http.StatusOK, []byte("/users/dcce1beb-aee6-4a4d-b724-94d470817323/password")},
		{"Invalid method", "Sqi1Z86I9632BndMkS632BNO4i1Z83v9KGd27FXs", http.MethodPost, http.StatusBadRequest, []byte("")},
		{"No token", "", http.MethodGet, http.StatusBadRequest, []byte("No token")},
		{"Expired token", "U5UC7FXsSq6I9632BndMkSO4i1Z83v9KGd2FDCDN", http.MethodGet, http.StatusBadRequest, []byte("Expired token")},
		{"Invalid token type/kind", "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2", http.MethodGet, http.StatusBadRequest, []byte("Invalid token")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := map[string]string{"token": tt.token}
			req, _ := http.NewRequest(tt.method, ts.URL+"/users/validateEmail", nil)
			q := req.URL.Query()
			for k, v := range query {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			code := rs.StatusCode
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if code == http.StatusOK {
				defer rs.Body.Close()
				respBody, err := ioutil.ReadAll(rs.Body)
				if err != nil {
					t.Fatal(err)
				}
				cookies := rs.Cookies()
				if len(cookies) == 0 {
					t.Errorf("Expect cookie 'sid'; got # %d cookies", len(cookies))
				} else {
					c := cookies[0]
					if c.Name != "sid" {
						t.Errorf("Expected a cookie names 'sid'; got %s", c.Name)
					}
				}
				if !bytes.Contains(respBody, tt.wantBody) {
					t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
				}
			}
		})
	}
}

func TestVerifyResetPw(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()
	tests := []struct {
		name     string
		token    string
		method   string
		wantCode int
		wantBody []byte
	}{
		{"Valid token", "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2", http.MethodGet, http.StatusOK, []byte("/users/dcce1beb-aee6-4a4d-b724-94d470817323/password")},
		{"Invalid method", "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2", http.MethodPost, http.StatusBadRequest, []byte("")},
		{"No token", "", http.MethodGet, http.StatusInternalServerError, []byte("No token")},
		{"Expired token", "U5UC7FXsSq6I9632BndMkSO4i1Z83v9KGd2FDCDN", http.MethodGet, http.StatusInternalServerError, []byte("Expired token")},
		{"Invalid token type/kind", "Sqi1Z86I9632BndMkS632BNO4i1Z83v9KGd27FXs", http.MethodGet, http.StatusInternalServerError, []byte("Invalid token")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := map[string]string{"token": tt.token}
			req, _ := http.NewRequest(tt.method, ts.URL+"/users/verifyResetPw", nil)
			q := req.URL.Query()
			for k, v := range query {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			code := rs.StatusCode
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if code == http.StatusOK {
				defer rs.Body.Close()
				respBody, err := ioutil.ReadAll(rs.Body)
				if err != nil {
					t.Fatal(err)
				}
				cookies := rs.Cookies()
				if len(cookies) == 0 {
					t.Errorf("Expect cookie 'sid'; got # %d cookies", len(cookies))
				} else {
					c := cookies[0]
					if c.Name != "sid" {
						t.Errorf("Expected a cookie names 'sid'; got %s", c.Name)
					}
				}
				if !bytes.Contains(respBody, tt.wantBody) {
					t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
				}
			}
		})
	}
}

// TestSetPasswordByClerkUser function tests user password update
func TestSetPasswordByClerkUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()
	tests := []struct {
		name         string
		UID          string
		userPassword string
		userConfirm  string
		method       string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", "validPa$$word", "validPa$$word", http.MethodPost, http.StatusOK, nil},
		{"UID of another user", "dcce1beb-aee6-4a4d-b724-94d470817323", "validPa$$word", "validPa$$word", http.MethodPost, http.StatusForbidden, []byte("Forbidden")},
		{"Wrong method", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", "validPa$$word", "validPa$$word", http.MethodGet, http.StatusMethodNotAllowed, []byte("Method Not Allowed")},
		{"Empty password", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", "", "", http.MethodPost, http.StatusBadRequest, []byte("Bad Request")},
		{"Short password", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", "pa$$wo", "pa$$wo", http.MethodPost, http.StatusBadRequest, []byte("Bad Request")},
		{"Invalid UID", "85c-4ff7-94ac-5afb5a1f0fcd", "pa$$wo432", "pa$$wo432", http.MethodPost, http.StatusBadRequest, []byte("Bad Request")},
		{"NonMatching Pws fields", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", "validPa$$word", "pa$$wo4329", http.MethodPost, http.StatusBadRequest, []byte("Bad Request")},
	}
	type postBody struct {
		Password string `json:"password"`
		Confirm  string `json:"confirm_password"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := &postBody{
				Password: tt.userPassword,
				Confirm:  tt.userConfirm}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}

			path := ts.URL + "/users/" + tt.UID + "/password"
			req, _ := http.NewRequest(tt.method, path, strings.NewReader(string(body)))
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e59d0bc-2fca-4aad-98d9-858709b66598",
			}
			req.AddCookie(cookie)
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			code := rs.StatusCode
			defer rs.Body.Close()
			respBody, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}
}

// TestLogout tests user logout and some middleware functions
func TestLogout(t *testing.T) {
	tests := []struct {
		name        string
		cookie      *http.Cookie
		wantCode    int
		wantBody    []byte
		wantCkValue string
	}{
		{
			name: "Valid session cookie",
			cookie: &http.Cookie{
				Name:  "sid",
				Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
			},
			wantCode:    http.StatusOK,
			wantBody:    []byte("User logged out"),
			wantCkValue: "",
		},
		{
			name: "Invalid session cookie",
			cookie: &http.Cookie{
				Name:  "sid",
				Value: "4976b9e8-a83f-40a5-b62e-a0ed206f953c",
			},
			wantCode:    http.StatusForbidden,
			wantBody:    []byte("Forbidden"),
			wantCkValue: "",
		},
		{
			name: "Expired session cookie",
			cookie: &http.Cookie{
				Name:  "sid",
				Value: "ef5de07d-4e78-4d8a-ab32-3529507a5072",
			},
			wantCode:    http.StatusForbidden,
			wantBody:    []byte("Forbidden"),
			wantCkValue: "",
		},
		{
			name: "Valid session cookie with an inactive user",
			cookie: &http.Cookie{
				Name:  "sid",
				Value: "03e77847-ff41-4be6-9fc8-166da95097b9",
			},
			wantCode:    http.StatusForbidden,
			wantBody:    []byte("Forbidden"),
			wantCkValue: "",
		},
		{
			name:        "No cookie",
			cookie:      nil,
			wantCode:    http.StatusForbidden,
			wantBody:    []byte("Forbidden"),
			wantCkValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := newTestApplication(t)
			ts := newTestServer(t, handlers.Routes())
			defer ts.Close()
			// prepare a post request with a valid session cookie
			req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users/logout", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			// send request
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer rs.Body.Close()
			respBody, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			if rs.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, rs.StatusCode)
			}

			if tt.wantCode == http.StatusOK {
				cookies := rs.Cookies()
				if len(cookies) == 0 {
					t.Errorf("Expect cookie 'sid'; got # %d cookies", len(cookies))
				} else {
					var ck *http.Cookie
					f := false
					for _, c := range cookies {
						if c.Name == "sid" {
							ck = c
							f = true
						}
					}
					if !f {
						t.Errorf("Expected a cookie named 'sid'; got none")
					}
					if f && ck.Value != tt.wantCkValue {
						t.Errorf("Expect cookie value: %s; got %s", tt.wantCkValue, ck.Value)
					}
					if f && ck.MaxAge != 0 {
						t.Errorf("Expect MaxAge == 0; got %d", ck.MaxAge)
					}
					if f && !bytes.Contains(respBody, tt.wantBody) {
						t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
					}
				}
			}
		})
	}
}

// TestLogin tests user login
func TestLogin(t *testing.T) {
	type postBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	tests := []struct {
		name      string
		wantName  []byte
		wantUID   []byte
		email     string
		wantPhone []byte
		pw        string
		body      *postBody
		method    string
		wantCode  int
	}{
		{
			name:      "Valid user",
			wantName:  []byte("Bob"),
			wantUID:   []byte("68b1d5e2-39dd-4713-8631-a08100383a0f"),
			email:     "bob@example.com",
			wantPhone: []byte("6544334535"),
			pw:        "test1234",
			body: &postBody{
				Email:    "bob@example.com",
				Password: "test1234",
			},
			method:   http.MethodPost,
			wantCode: http.StatusOK,
		},
		{
			name:      "Bad Method",
			wantName:  []byte(""),
			wantUID:   []byte(""),
			email:     "bob@example.com",
			wantPhone: []byte(""),
			pw:        "test1234",
			body: &postBody{
				Email:    "bob@example.com",
				Password: "test1234",
			},
			method:   http.MethodGet,
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "Empty password",
			wantName:  []byte(""),
			wantUID:   []byte(""),
			email:     "bob@example.com",
			wantPhone: []byte(""),
			pw:        "",
			body: &postBody{
				Email:    "bob@example.com",
				Password: "",
			},
			method:   http.MethodPost,
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "Bad post body",
			wantName:  []byte(""),
			wantUID:   []byte(""),
			email:     "bob@example.com",
			wantPhone: []byte(""),
			pw:        "test1234",
			body:      &postBody{},
			method:    http.MethodPost,
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "Bad email",
			wantName:  []byte(""),
			wantUID:   []byte(""),
			email:     "bob",
			wantPhone: []byte(""),
			pw:        "24512341",
			body: &postBody{
				Email:    "bob",
				Password: "test1234",
			},
			method:   http.MethodPost,
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "Invalid user",
			wantName:  []byte("Invalid"),
			wantUID:   []byte(""),
			email:     "mary@example.com",
			wantPhone: []byte("232222"),
			pw:        "anotherpw1234",
			body: &postBody{
				Email:    "mary@example.com",
				Password: "anotherpw1234",
			},
			method:   http.MethodPost,
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := newTestApplication(t)
			ts := newTestServer(t, handlers.Routes())
			defer ts.Close()
			reqBody := &tt.body
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			req, _ := http.NewRequest(tt.method, ts.URL+"/users/login", strings.NewReader(string(body)))
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			code := rs.StatusCode
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			// extras checks if login was successful
			if code == http.StatusOK {
				defer rs.Body.Close()
				respBody, err := ioutil.ReadAll(rs.Body)
				if err != nil {
					t.Fatal(err)
				}
				cookies := rs.Cookies()
				if len(cookies) == 0 {
					t.Errorf("Expect cookie 'sid'; got # %d cookies", len(cookies))
				} else {
					c := cookies[0]
					if c.Name != "sid" {
						t.Errorf("Expected a cookie names 'sid'; got %s", c.Name)
					}
				}
				if !bytes.Contains(respBody, tt.wantUID) {
					t.Errorf("want body %s to contain %q", respBody, tt.wantUID)
				}
				if !bytes.Contains(respBody, tt.wantName) {
					t.Errorf("want body %s to contain %q", respBody, tt.wantName)
				}
				if !bytes.Contains(respBody, tt.wantPhone) {
					t.Errorf("want body %s to contain %q", respBody, tt.wantPhone)
				}
				if !bytes.Contains(respBody, []byte(tt.email)) {
					t.Errorf("want body %s to contain %q", respBody, tt.email)
				}
			}
		})
	}
}

// TestGetUserByEmail function tests retrieving user data by email
func TestGetUserByEmailByAdmin(t *testing.T) {
	handlers := newTestApplication(t)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	tests := []struct {
		name     string
		email    string
		wantCode int
		wantBody []byte
	}{
		{"User not found", "test@email.com", http.StatusNotFound, []byte("Not Found")},
		{"Bad user request", "test.com", http.StatusBadRequest, []byte("Bad Request")},
		{"Empty request", "", http.StatusBadRequest, []byte("Bad Request")},
		{"Valid user", "bob@example.com", http.StatusOK, []byte("bob@example.com")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := map[string]string{"email": tt.email}
			req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
			q := req.URL.Query()
			for k, v := range query {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
			}
			req.AddCookie(cookie)
			code, _, respBody := ts.getQueryReq(t, req)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}

	// t.Run("Reject POST request", func(t *testing.T) {
	// 	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/getUserByEmail", nil)
	// 	q := req.URL.Query()
	// 	q.Add("email", "bob@example.com")
	// 	req.URL.RawQuery = q.Encode()
	// 	cookie := &http.Cookie{
	// 		Name:  "sid",
	// 		Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
	// 	}
	// 	req.AddCookie(cookie)
	// 	rs, err := ts.Client().Do(req)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	defer rs.Body.Close()
	// 	if rs.StatusCode != http.StatusBadRequest {
	// 		t.Errorf("want %d; got %d", http.StatusBadRequest, rs.StatusCode)
	// 	}
	// })
}

func TestGetUserByEmailByClerkUser(t *testing.T) {
	handlers := newTestApplication(t)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	tests := []struct {
		name     string
		email    string
		wantCode int
		wantBody []byte
	}{
		{"User not found", "test@email.com", http.StatusForbidden, []byte("Forbidden")},
		{"Bad user request", "test.com", http.StatusForbidden, []byte("Forbidden")},
		{"Empty request", "", http.StatusForbidden, []byte("Forbidden")},
		// normal users are not supposed to use this endpoint, they already receive
		// their own relevant data when login
		{"Valid user", "valid@user.com", http.StatusForbidden, []byte("Forbidden")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := map[string]string{"email": tt.email}
			req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
			q := req.URL.Query()
			for k, v := range query {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e59d0bc-2fca-4aad-98d9-858709b66598",
			}
			req.AddCookie(cookie)
			code, _, respBody := ts.getQueryReq(t, req)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}

	// t.Run("Reject POST request", func(t *testing.T) {
	// 	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/getUserByEmail", nil)
	// 	q := req.URL.Query()
	// 	q.Add("email", "valid@user.com")
	// 	req.URL.RawQuery = q.Encode()
	// 	cookie := &http.Cookie{
	// 		Name:  "sid",
	// 		Value: "4e59d0bc-2fca-4aad-98d9-858709b66598",
	// 	}
	// 	req.AddCookie(cookie)
	// 	rs, err := ts.Client().Do(req)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	defer rs.Body.Close()
	// 	if rs.StatusCode != http.StatusForbidden {
	// 		t.Errorf("want %d; got %d", http.StatusBadRequest, rs.StatusCode)
	// 	}
	// })
}

// TestGetUserName function test end point for getting a user name
func TestGetUserNameByClerkUser(t *testing.T) {
	handlers := newTestApplication(t)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	tests := []struct {
		name     string
		uID      string
		wantCode int
		wantBody []byte
	}{
		{"Valid user", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", http.StatusOK, []byte("Valid User")},
		{"Another user's data", "68b1d5e2-39dd-4713-8631-a08100383a0f", http.StatusForbidden, []byte("Forbidden")},
	}
	type getBody struct {
		Email string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := ts.URL + "/users/" + tt.uID + "/name"
			t.Log(url)
			req, _ := http.NewRequest(http.MethodGet, url, nil)

			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e59d0bc-2fca-4aad-98d9-858709b66598",
			}
			req.AddCookie(cookie)
			code, _, respBody := ts.getQueryReq(t, req)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}
}

func TestGetUserNameByAdminUser(t *testing.T) {
	handlers := newTestApplication(t)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	// Bob is admin
	tests := []struct {
		name     string
		uID      string
		wantCode int
		wantBody []byte
	}{
		{"Another user's data", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", http.StatusOK, []byte("Valid User")},
		{"Same user", "68b1d5e2-39dd-4713-8631-a08100383a0f", http.StatusOK, []byte("Bob")},
	}
	type getBody struct {
		Email string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := ts.URL + "/users/" + tt.uID + "/name"
			t.Log(url)
			req, _ := http.NewRequest(http.MethodGet, url, nil)

			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
			}
			req.AddCookie(cookie)
			code, _, respBody := ts.getQueryReq(t, req)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}
}

func TestGetUserEmailByClerkUser(t *testing.T) {
	handlers := newTestApplication(t)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	tests := []struct {
		name     string
		uID      string
		method   string
		wantCode int
		wantBody []byte
	}{
		{"Valid user", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", http.MethodGet, http.StatusOK, []byte("valid@user.com")},
		{"Another user's data", "68b1d5e2-39dd-4713-8631-a08100383a0f", http.MethodGet, http.StatusForbidden, []byte("Forbidden")},
		{"Wrong method", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", http.MethodPost, http.StatusMethodNotAllowed, []byte("Method Not Allowed")},
	}
	type getBody struct {
		Email string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := ts.URL + "/users/" + tt.uID + "/email"
			t.Log(url)
			req, _ := http.NewRequest(tt.method, url, nil)

			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e59d0bc-2fca-4aad-98d9-858709b66598",
			}
			req.AddCookie(cookie)
			code, _, respBody := ts.getQueryReq(t, req)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}
}

func TestGetUserEmailByAdminUser(t *testing.T) {
	handlers := newTestApplication(t)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	// Bob is admin
	tests := []struct {
		name     string
		uID      string
		method   string
		wantCode int
		wantBody []byte
	}{
		{"Another user's data", "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd", http.MethodGet, http.StatusOK, []byte("valid@user.com")},
		{"Same user", "68b1d5e2-39dd-4713-8631-a08100383a0f", http.MethodGet, http.StatusOK, []byte("bob@example.com")},
		{"Wrong method", "68b1d5e2-39dd-4713-8631-a08100383a0f", http.MethodPut, http.StatusMethodNotAllowed, []byte("Method Not Allowed")},
	}
	type getBody struct {
		Email string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := ts.URL + "/users/" + tt.uID + "/email"
			t.Log(url)
			req, _ := http.NewRequest(tt.method, url, nil)

			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
			}
			req.AddCookie(cookie)
			code, _, respBody := ts.getQueryReq(t, req)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}
}

// TestAddUser in /signup end point function tests user insertion using mocked depedencies (end-to-end)
// func TestAddUser(t *testing.T) {
// 	handlers := newTestApplication(t)
// 	ts := newTestServer(t, handlers.Routes())
// 	defer ts.Close()
// 	tests := []struct {
// 		name         string
// 		userName     string
// 		userEmail    string
// 		userPassword string
// 		wantCode     int
// 		wantBody     []byte
// 	}{
// 		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", http.StatusOK, []byte("bob@example.com")},
// 		{"Empty name", "", "bob@example.com", "validPa$$word", http.StatusBadRequest, []byte("This field cannot be blank")},
// 		{"Empty email", "Bob", "", "validPa$$word", http.StatusBadRequest, []byte("This field cannot be blank")},
// 		{"Empty password", "Bob", "bob@example.com", "", http.StatusBadRequest, []byte("This field cannot be blank")},
// 		{"Short password", "Bob", "bob@example.com", "pa$$wo", http.StatusBadRequest, []byte("This field is too short (minimum is 8 characters)")},
// 		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", http.StatusBadRequest, []byte("address already in use")},
// 		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", http.StatusBadRequest, []byte("This field is invalid")},
// 		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", http.StatusBadRequest, []byte("This field is invalid")},
// 		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", http.StatusBadRequest, []byte("This field is invalid")},
// 	}

// 	type postBody struct {
// 		Name, Email, Password string
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			reqBody := &postBody{Name: tt.userName,
// 				Email:    tt.userEmail,
// 				Password: tt.userPassword}
// 			body, err := json.Marshal(reqBody)
// 			if err != nil {
// 				t.Log(err)
// 			}
// 			code, _, respBody := ts.post(t, "/signup", "applicaton/json", strings.NewReader(string(body)))
// 			if code != tt.wantCode {
// 				t.Errorf("want %d; got %d", tt.wantCode, code)
// 			}
// 			if !bytes.Contains(respBody, tt.wantBody) {
// 				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
// 			}
// 		})
// 	}
// 	t.Run("Empty req body", func(t *testing.T) {
// 		code, _, _ := ts.post(t, "/signup", "applicaton/json", strings.NewReader(""))
// 		if code != http.StatusInternalServerError {
// 			t.Errorf("want body %d to contain %d", http.StatusInternalServerError, code)
// 		}
// 	})
// 	t.Run("Bad req content type", func(t *testing.T) {
// 		code, _, _ := ts.post(t, "/signup", "applicaton/text", strings.NewReader(""))
// 		if code != http.StatusInternalServerError {
// 			t.Errorf("want body %d to contain %d", http.StatusInternalServerError, code)
// 		}
// 	})
// 	t.Run("Bad body", func(t *testing.T) {
// 		mockedReader := &BadReader{}
// 		code, _, _ := ts.post(t, "/signup", "applicaton/json", mockedReader)
// 		// t.Log(code)
// 		// t.Log(string(respBody))
// 		if code != http.StatusInternalServerError {
// 			t.Errorf("want body %d to contain %d", http.StatusInternalServerError, code)
// 		}
// 	})

// 	t.Run("Reject GET request", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/signup", strings.NewReader(""))
// 		req.Header.Set("Content-Type", "applicaton/json")
// 		rs, err := ts.Client().Do(req)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		defer rs.Body.Close()
// 		if rs.StatusCode != http.StatusBadRequest {
// 			t.Errorf("want %d; got %d", http.StatusBadRequest, rs.StatusCode)
// 		}
// 	})

// }

// type BadReader struct{}

// func (d *BadReader) Read(p []byte) (i int, e error) {
// 	return 0, errors.New("BadBody")
// }

// TestAddUserByAdmin function tests user insertion using mocked depedencies (end-to-end)
func TestAddUserByAdmin(t *testing.T) {
	// use a mail recorder
	// https://tmichel.github.io/2014/10/12/golang-send-test-email/
	rec := new(mocks.EmailRecorder)
	mailer := mocks.NewMailerMock(rec)
	handlers := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		mocks.NewUserSvc(),
		mailer,
		mocks.NewTokenMockSvc())

	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		userPhone    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", "12345", http.StatusOK, []byte("bob@example.com")},
		{"Empty name", "", "bob@example.com", "validPa$$word", "12345", http.StatusBadRequest, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", "12345", http.StatusBadRequest, []byte("This field cannot be blank")},
		// {"Empty password", "Bob", "bob@example.com", "", "12345", http.StatusBadRequest, []byte("This field cannot be blank")},
		// {"Short password", "Bob", "bob@example.com", "pa$$wo", "12345", http.StatusBadRequest, []byte("This field is too short (minimum is 8 characters)")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", "12345", http.StatusBadRequest, []byte("address already in use")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", "12345", http.StatusBadRequest, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", "12345", http.StatusBadRequest, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", "12345", http.StatusBadRequest, []byte("This field is invalid")},
	}

	type postBody struct {
		Name, Email, Password, Phone string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := &postBody{Name: tt.userName,
				Email:    tt.userEmail,
				Phone:    tt.userPhone,
				Password: tt.userPassword,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", strings.NewReader(string(body)))
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
			}
			req.AddCookie(cookie)
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			code := rs.StatusCode
			defer rs.Body.Close()
			respBody, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			if code == http.StatusOK {
				bMsg := []byte(rec.Msg)
				if !(rec.Email == tt.userEmail) ||
					!(rec.Name == tt.userName) ||
					!bytes.Contains(bMsg, []byte(tt.userName)) ||
					!bytes.Contains(bMsg, []byte("users/validateEmail?token=")) {
					t.Errorf("missing data from email msg sent to user")
				}
			}
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(respBody, tt.wantBody) {
				t.Errorf("want body %s to contain %q", respBody, tt.wantBody)
			}
		})
	}
	t.Run("Empty req body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", strings.NewReader(""))
		cookie := &http.Cookie{
			Name:  "sid",
			Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
		}
		req.AddCookie(cookie)
		rs, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		code := rs.StatusCode
		if code != http.StatusInternalServerError {
			t.Errorf("want body %d to contain %d", http.StatusInternalServerError, code)
		}
	})
	// t.Run("Bad req content type", func(t *testing.T) {
	// 	code, _, _ := ts.post(t, "/users", "applicaton/text", strings.NewReader(""))
	// 	if code != http.StatusInternalServerError {
	// 		t.Errorf("want body %d to contain %d", http.StatusInternalServerError, code)
	// 	}
	// })
	// t.Run("Bad body", func(t *testing.T) {
	// 	mockedReader := &BadReader{}
	// 	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", mockedReader)
	// 	cookie := &http.Cookie{
	// 		Name:  "sid",
	// 		Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
	// 	}
	// 	req.AddCookie(cookie)
	// 	rs, err := ts.Client().Do(req)
	// 	if err != nil {
	// 		t.Log(err)
	// 		t.Fatal(err)
	// 	}
	// 	code := rs.StatusCode
	// 	if code != http.StatusInternalServerError {
	// 		t.Errorf("want body %d to contain %d", http.StatusInternalServerError, code)
	// 	}
	// })

	// t.Run("Reject GET request", func(t *testing.T) {
	// 	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", strings.NewReader(""))
	// 	req.Header.Set("Content-Type", "applicaton/json")
	// 	rs, err := ts.Client().Do(req)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	defer rs.Body.Close()
	// 	if rs.StatusCode != http.StatusBadRequest {
	// 		t.Errorf("want %d; got %d", http.StatusBadRequest, rs.StatusCode)
	// 	}
	// })

}

func TestDefaultUsersHandlerCase(t *testing.T) {
	handlers := newTestApplication(t)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	testCases := []struct {
		desc     string
		path     string
		wantCode int
	}{
		{
			desc:     "Bad path",
			path:     "/users/85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd/badpath",
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, ts.URL+tC.path, nil)
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
			}
			req.AddCookie(cookie)
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			code := rs.StatusCode
			if code != tC.wantCode {
				t.Errorf("want %d; got %d", tC.wantCode, code)
			}

		})
	}
}

func TestUsersNoPathBadMethods(t *testing.T) {
	handlers := newTestApplication(t)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	testCases := []struct {
		desc     string
		method   string
		wantCode int
	}{
		{
			desc:     "Patch",
			method:   http.MethodPatch,
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Put",
			method:   http.MethodPut,
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Delete",
			method:   http.MethodDelete,
			wantCode: http.StatusMethodNotAllowed,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest(tC.method, ts.URL+"/users", nil)
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
			}
			req.AddCookie(cookie)
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			code := rs.StatusCode
			if code != tC.wantCode {
				t.Errorf("want %d; got %d", tC.wantCode, code)
			}

		})
	}
}

// TestPasswordReset function tests user insertion using mocked depedencies (end-to-end)
func TestPasswordResetRequest(t *testing.T) {
	// use a mail recorder
	// https://tmichel.github.io/2014/10/12/golang-send-test-email/
	rec := new(mocks.EmailRecorder)
	mailer := mocks.NewMailerMock(rec)
	app := web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		mocks.NewUserSvc(),
		mailer,
		mocks.NewTokenMockSvc())

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	tests := []struct {
		name      string
		userName  string
		userEmail string
		wantCode  int
	}{
		{"Valid submission", "Bob", "bob@example.com", http.StatusOK},
		{"Disabled account", "Barack Obama", "bobama@somewhere.com", http.StatusBadRequest},
		{"Empty email", "Bob", "", http.StatusBadRequest},
		{"Non-existing user", "Mark", "Mark@somewhere.com", http.StatusOK},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", http.StatusBadRequest},
		{"Invalid email (missing @)", "Bob", "bobexample.com", http.StatusBadRequest},
		{"Invalid email (missing local part)", "Bob", "@example.com", http.StatusBadRequest},
	}
	type postBody struct {
		Email string `json:"email"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := &postBody{
				Email: tt.userEmail,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Log(err)
			}
			req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users/resetPassword", strings.NewReader(string(body)))

			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			code := rs.StatusCode
			defer rs.Body.Close()
			respBody, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(respBody))
			if code == http.StatusOK {
				bMsg := []byte(rec.Msg)
				if tt.name == "Non-existing user" {
					if !(rec.Email == tt.userEmail) {
						t.Errorf("missing data from email msg sent to user")
					}
				} else {
					if !(rec.Email == tt.userEmail) ||
						!(rec.Name == tt.userName) ||
						!bytes.Contains(bMsg, []byte(tt.userName)) ||
						!bytes.Contains(bMsg, []byte("users/verifyResetPw?token=")) {
						t.Errorf("missing data from email msg sent to user")
					}
				}
			}
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
		})
	}
}

// TestGetATestGetAllUsersByAdmin
func TestGetAllUsersByAdmin(t *testing.T) {
	// use a mail recorder
	// https://tmichel.github.io/2014/10/12/golang-send-test-email/
	rec := new(mocks.EmailRecorder)
	mailer := mocks.NewMailerMock(rec)
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
		mailer,
		mocks.NewTokenMockSvc())

	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
	tests := []struct {
		name         string
		previous     string
		next         string
		pgSize       string
		wantSize     int
		wantCode     int
		wantResponse string
	}{
		{"Valid submission", "", "", "6", 6, http.StatusOK, "bob@example.com"},
		{"Invalid submission", "dmFsaWRAdXNlci5jb20=", "dGJsZWVAc29tZXdoZXJlLmNvbQ==", "6", 6, http.StatusBadRequest, ""},
	}

	// page encapsulates data for pagination
	type page struct {
		StartCursor     string `json:"startCursor"`
		HasPreviousPage bool   `json:"hasprevious"`
		EndCursor       string `json:"endCursor"`
		HasNextPage     bool   `json:"hasnext"`
		Users           []user.User
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//t.Logf("previous: %s\n next: %s\n", tt.previous, tt.next)
			query := map[string]string{
				"previous": tt.previous,
				"next":     tt.next,
				"pgSize":   tt.pgSize,
			}
			req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
			q := req.URL.Query()
			for k, v := range query {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			cookie := &http.Cookie{
				Name:  "sid",
				Value: "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
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
			var p page
			if tt.wantCode == http.StatusOK {
				defer rs.Body.Close()
				err = json.NewDecoder(rs.Body).Decode(&p)
				if err != nil {
					t.Error("bad response body")
				}
			}

			//t.Logf("page received: %v", p)
			if tt.wantResponse != "" {
				contains := false
				for _, u := range p.Users {
					if u.Email == tt.wantResponse {
						contains = true
						break
					}
				}
				if !contains {
					t.Errorf("want response to contain %v", tt.wantResponse)
				}
			}
		})
	}

}
