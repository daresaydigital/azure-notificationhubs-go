package utils_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/daresaydigital/azure-notificationhubs-go/utils"
)

const (
	errfmt = "Expected %s: \n%v\ngot:\n%v"
)

func Test_WithoutOnRequestHook(t *testing.T) {
	var (
		expectedError = "Get http://0.0.0.0:9999"
		client        = NewHubHTTPClient()
		req, _        = http.NewRequest("GET", "http://0.0.0.0:9999", nil)
	)
	client.SetTimeout(1 * time.Millisecond)
	_, resp, err := client.Exec(req)
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf(errfmt, "error from request", expectedError, err)
	}
	if resp != nil {
		t.Errorf(errfmt, "no response back from dummy request", "nil", err)
	}
}

func Test_OnRequestHook(t *testing.T) {
	var (
		expectedError = "Get http://0.0.0.0:9999"
		client        = NewHubHTTPClient()
		req, _        = http.NewRequest("GET", "http://0.0.0.0:9999", nil)
	)
	client.SetTimeout(1 * time.Millisecond)
	called := false
	client.OnRequest(func(r *http.Request) {
		if r != req {
			t.Errorf(errfmt, "should be the same request", req, r)
		}
		called = true
	})
	_, resp, err := client.Exec(req)
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf(errfmt, "error from request", expectedError, err)
	}
	if resp != nil {
		t.Errorf(errfmt, "no response back from dummy request", "nil", err)
	}
	if !called {
		t.Errorf(errfmt, "on request hook should be called", true, called)
	}
}
