package web_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/web"
)

func TestAddCookie(t *testing.T) {
	ss := make(map[string]http.SameSite)
	ss["lax"] = http.SameSiteLaxMode
	ss["strict"] = http.SameSiteStrictMode
	ss["default"] = http.SameSiteDefaultMode

	testCases := []struct {
		desc     string
		sameSite string
	}{
		{
			desc:     "Lax Cookie",
			sameSite: "lax",
		},
		{
			desc:     "Strict Cookie",
			sameSite: "strict",
		},
		{
			desc:     "Default Cookie",
			sameSite: "anyValue",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			app := web.New(
				log.New(ioutil.Discard, "", 0),
				log.New(ioutil.Discard, "", 0),
				&web.CkProps{
					Name:     "sid",
					HTTPOnly: false,
					Secure:   false,
					SameSite: tC.sameSite,
				},
				mocks.NewSessionSvc(),
				mocks.NewUserSvc(),
				mocks.NewMailerMock(nil),
				mocks.NewTokenMockSvc(),
				nil, nil, nil)
			rr := httptest.NewRecorder()
			app.AddCookie(rr, "1")
			rs := rr.Result()
			cookies := rs.Cookies()
			if len(cookies) == 0 {
				t.Errorf("Expect cookie 'sid'; got # %d cookies", len(cookies))
			} else {
				c := cookies[0]
				if c.Name != "sid" {
					t.Errorf("Expected a cookie names 'sid'; got %s", c.Name)
				}
				if c.SameSite != ss[tC.sameSite] && c.SameSite != http.SameSiteDefaultMode {
					t.Errorf("Expected cookie SameSite %v; got %v", tC.sameSite, c.SameSite)
				}

			}
		})
	}
}
