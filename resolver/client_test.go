package resolver

import (
	"bytes"
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

var requestedUrl string

type fakeHttpCli struct {
	statusCode      int
	isErrorResponse bool
	error           error
	responseStr     string
	bodyError       bool
}

func (faker *fakeHttpCli) Do(req *http.Request) (*http.Response, error) {
	requestedUrl = req.URL.Path
	rdr := bytes.NewReader([]byte(faker.responseStr))
	if faker.isErrorResponse {
		fmt.Println("1")
		return nil, faker.error
	}
	if faker.bodyError {
		return &http.Response{Body: nil}, nil
	}
	return &http.Response{Body: ioutil.NopCloser(rdr), StatusCode: faker.statusCode}, nil
}

func TestResolverGet(t *testing.T) {
	url := "/some/resource"
	responseStr := "Hello Kitty"
	Convey("Client returns a 200 status code and the expected bytes for a successful response.", t, func() {
		Client = &fakeHttpCli{statusCode: 200, responseStr: responseStr}
		b, err := Get(url)
		So(err, ShouldBeNil)
		So(string(b), ShouldEqual, responseStr)
		So(requestedUrl, ShouldEqual, url)

	})

	Convey("Client returns error if get request fails.", t, func() {
		expectedErr := errors.New("Somthing went wrong")
		Client = &fakeHttpCli{statusCode: 500, isErrorResponse: true, error: expectedErr}
		b, err := Get("/")
		So(err, ShouldEqual, expectedErr)
		So(string(b), ShouldEqual, "")
	})

	Convey("Client returns error if reading the response.body fails.", t, func() {
		expectedErr := errors.New("Error reading body")
		readResponseBody = func(r io.Reader) ([]byte, error) {
			return make([]byte, 0), errors.New("Error reading body")
		}
		Client = &fakeHttpCli{statusCode: 500, bodyError: true, error: expectedErr}
		b, err := Get("/")
		So(err.Error(), ShouldEqual, expectedErr.Error())
		So(string(b), ShouldEqual, "")
	})

	Convey("Client returns empty bytes slice and error is response status is not 200", t, func() {
		expectedErr := errors.New("Response status code is not 200")
		readResponseBody = ioutil.ReadAll
		Client = &fakeHttpCli{statusCode: 500, isErrorResponse: false, error: expectedErr}
		b, err := Get("/")
		So(b, ShouldBeEmpty)
		So(err.Error(), ShouldEqual, expectedErr.Error())
	})
}
