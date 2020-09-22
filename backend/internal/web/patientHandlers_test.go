package web_test

import (
	"io/ioutil"
	"log"
	"testing"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/web"
)

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
		mocks.NewUserSvc(),
		nil,
		mocks.NewTokenMockSvc(),
	)
	ts := newTestServer(t, handlers.Routes())
	defer ts.Close()
}
