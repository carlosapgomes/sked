package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/user"
	"carlosapgomes.com/sked/internal/web"
)

func TestIsSameUserOrAdmin(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		ctxData      map[string]interface{}
		userData     user.User
		wantResponse bool
	}{
		{
			name:   "Same User GET",
			method: http.MethodGet,
			ctxData: map[string]interface{}{
				"isAuthenticaded": true,
				"user": user.User{
					ID:    "adea27f8-4091-4908-9ad9-3d68e198b488",
					Email: "user@test.com",
				},
			},
			userData: user.User{
				ID:    "adea27f8-4091-4908-9ad9-3d68e198b488",
				Email: "user@test.com",
			},
			wantResponse: true,
		},
		{
			name:   "Different User GET",
			method: http.MethodGet,
			ctxData: map[string]interface{}{
				"isAuthenticaded": true,
				"user": user.User{
					ID:    "adea27f8-4091-4908-9ad9-3d68e198b488",
					Email: "user@test.com",
				},
			},
			userData: user.User{
				ID:    "b9b127c8-6deb-4d16-88ff-cca15bf9ff47",
				Email: "anotheruser@test.com",
			},
			wantResponse: false,
		}, {
			name:   "Different User GET, Admin",
			method: http.MethodGet,
			ctxData: map[string]interface{}{
				"isAuthenticaded": true,
				"user": user.User{
					ID:    "adea27f8-4091-4908-9ad9-3d68e198b488",
					Email: "user@test.com",
					Roles: []string{user.RoleCommon, user.RoleAdmin},
				},
			},
			userData: user.User{
				ID:    "b9b127c8-6deb-4d16-88ff-cca15bf9ff47",
				Email: "anotheruser@test.com",
			},
			wantResponse: true,
		},
		{
			name:   "Same User POST",
			method: http.MethodPost,
			ctxData: map[string]interface{}{
				"isAuthenticaded": true,
				"user": user.User{
					ID:    "adea27f8-4091-4908-9ad9-3d68e198b488",
					Email: "user@test.com",
				},
			},
			userData: user.User{
				ID:    "adea27f8-4091-4908-9ad9-3d68e198b488",
				Email: "user@test.com",
			},
			wantResponse: true,
		},
		{
			name:   "Different User POST",
			method: http.MethodPost,
			ctxData: map[string]interface{}{
				"isAuthenticaded": true,
				"user": user.User{
					ID:    "adea27f8-4091-4908-9ad9-3d68e198b488",
					Email: "user@test.com",
				},
			},
			userData: user.User{
				ID:    "b9b127c8-6deb-4d16-88ff-cca15bf9ff47",
				Email: "anotheruser@test.com",
			},
			wantResponse: false,
		},
		{
			name:   "Different User POST, Admin",
			method: http.MethodPost,
			ctxData: map[string]interface{}{
				"isAuthenticaded": true,
				"user": user.User{
					ID:    "adea27f8-4091-4908-9ad9-3d68e198b488",
					Email: "user@test.com",
					Roles: []string{user.RoleCommon, user.RoleAdmin},
				},
			},
			userData: user.User{
				ID:    "b9b127c8-6deb-4d16-88ff-cca15bf9ff47",
				Email: "anotheruser@test.com",
			},
			wantResponse: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := web.New(log.New(ioutil.Discard, "", 0),
				log.New(ioutil.Discard, "", 0),
				nil,
				mocks.NewSessionSvc(),
				mocks.NewUserSvc(),
				mocks.NewMailerMock(nil),
				mocks.NewTokenMockSvc())
			var req *http.Request
			if tt.method == http.MethodGet {
				req, _ = http.NewRequest(tt.method, "/", nil)
				q := req.URL.Query()
				q.Add("email", tt.userData.Email)
				req.URL.RawQuery = q.Encode()
			} else {
				body, err := json.Marshal(map[string]string{
					"email": tt.userData.Email,
				})
				if err != nil {
					t.Error("could not create request body")
				}
				req, err = http.NewRequest(tt.method, "/", bytes.NewBuffer(body))
				if err != nil {
					t.Error("could not create request")
				}
				req.Header.Set("Content-type", "application/json")
			}
			usr := tt.ctxData["user"].(user.User)
			ctx := context.WithValue(req.Context(), web.ContextKeyIsAuthenticated, tt.ctxData["isAuthenticated"])
			ctx = context.WithValue(ctx, web.ContextKeyUser, &usr)
			resp := app.IsSameUserOrAdmin(req.WithContext(ctx))
			if resp != tt.wantResponse {
				t.Errorf("want response %t; got %t", tt.wantResponse, resp)
			}
		})
	}
}
