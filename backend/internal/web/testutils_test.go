package web_test

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/web"
)

// MailerMock holds an instance of a mailer
var MailerMock *mocks.MailerMockSvc

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *web.App {
	return web.New(
		log.New(ioutil.Discard, "", 0),
		log.New(ioutil.Discard, "", 0),
		&web.CkProps{
			Name:     "sid",
			HTTPOnly: false,
			Secure:   false,
		},
		mocks.NewSessionSvc(),
		mocks.NewUserSvc(),
		mocks.NewMailerMock(nil),
		mocks.NewTokenMockSvc(),
		nil, nil, nil)
}

// Define a custom testServer type which anonymously embeds a httptest.Server
// instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
// We then use the httptest.NewServer() function to create a new test
// server, passing in the value returned by our app.routes() method as the
// handler for the server. This starts up a HTTP server which listens on a
// randomly-chosen port of your local machine for the duration of the test.

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add the cookie jar to the client, so that response cookies are stored
	// and then sent with subsequent requests.
	ts.Client().Jar = jar

	// Disable redirect-following for the client. Essentially this function
	// is called after a 3xx response is received by the client, and returning
	// the http.ErrUseLastResponse error forces it to immediately return the
	// received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	// The network address that the test server is listening on is contained
	// in the ts.URL field. We can use this along with the ts.Client().Get()
	// method to make a GET /ping request against the test server. This
	// returns a http.Response struct containing the response.
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}
func (ts *testServer) getQuery(t *testing.T, urlPath string, query map[string]string) (int, http.Header, []byte) {
	// The network address that the test server is listening on is contained
	// in the ts.URL field. We can use this along with the ts.Client().Get()
	// method to make a GET /ping request against the test server. This
	// returns a http.Response struct containing the response.
	// To make a request with custom headers, use NewRequest and Client.Do.

	req, _ := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}

// same as above but allow a pre-made request
func (ts *testServer) getQueryReq(t *testing.T, r *http.Request) (int, http.Header, []byte) {
	// The network address that the test server is listening on is contained
	// in the ts.URL field. We can use this along with the ts.Client().Get()
	// method to make a GET /ping request against the test server. This
	// returns a http.Response struct containing the response.
	// To make a request with custom headers, use NewRequest and Client.Do.

	// req, _ := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	// q := req.URL.Query()
	// for k, v := range query {
	// 	q.Add(k, v)
	// }
	// req.URL.RawQuery = q.Encode()
	rs, err := ts.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}

// Implement a Post method on our custom testServer type. This makes a POST
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) post(t *testing.T, urlPath, contentType string, body io.Reader) (int, http.Header, []byte) {
	// The network address that the test server is listening on is contained
	// in the ts.URL field. We can use this along with the ts.Client().Post()
	// method to make a post request
	// t.Log(ts.URL)
	// t.Log(urlPath)
	rs, err := ts.Client().Post(ts.URL+urlPath, contentType, body)
	if err != nil {
		if strings.Contains(err.Error(), errors.New("BadBody").Error()) {
			return 500, nil, nil
		} else {
			t.Fatal(err)
		}
	}
	defer rs.Body.Close()
	respBody, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, respBody
}
