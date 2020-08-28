package services_test

import (
	"net/http"
	"os"
	"testing"

	"carlosapgomes.com/sked/internal/services"
)

func Test(t *testing.T) {
	var sgKey string
	if key, ok := os.LookupEnv("SENDGRID_API_KEY"); ok {
		sgKey = key
	} else {
		t.Log("skipping")
		t.Skip("mailer: missing Sendgrid key from environment")
	}
	testCases := []struct {
		desc      string
		to        string
		from      string
		msg       string
		wantError error
		wantCode  int
	}{
		{
			desc:      "Valid email",
			to:        "capgomes2015+sg@gmail.com",
			from:      "sked.manager@gmail.com",
			msg:       "<h3>Hello</h3>",
			wantError: nil,
			wantCode:  http.StatusAccepted,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mailer := services.NewMailerService(sgKey, "sked.manager", tC.from)
			r, err := mailer.Send("Carlos", tC.to, "ts sg", tC.msg)
			if err != tC.wantError {
				t.Errorf("want \n%v\n; got \n%v\n", tC.wantError, err)
			}
			if r.Code != tC.wantCode {
				t.Errorf("want \n%v\n; got \n%v\n", tC.wantCode, r.Code)
				t.Errorf(r.Msg)
			}
		})
	}
}
