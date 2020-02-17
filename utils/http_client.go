// Package utils are the the place for replaceable utils in the module
package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	// HTTPClient interface, replaceable for testing or custom implementation
	HTTPClient interface {
		Exec(*http.Request) ([]byte, *http.Response, error)
		OnRequest(OnRequestFunc)
		SetTimeout(time.Duration)
	}

	hubHTTPClient struct {
		httpClient *http.Client
		onRequest  OnRequestFunc
	}

	// OnRequestFunc takes a request and returns nothing, add it using the OnRequest hook
	OnRequestFunc func(*http.Request)
)

// NewHubHTTPClient is creating the default client
func NewHubHTTPClient() HTTPClient {
	return &hubHTTPClient{
		httpClient: &http.Client{},
	}
}

// Exec executes notification hub http request and handles the response
func (hc *hubHTTPClient) Exec(req *http.Request) ([]byte, *http.Response, error) {
	if hc.onRequest != nil {
		hc.onRequest(req)
	}
	return handleResponse(hc.httpClient.Do(req))
}

// OnRequest adds an optional hook to add more logging or other upon a request from the hub
func (hc *hubHTTPClient) OnRequest(fun OnRequestFunc) {
	hc.onRequest = fun
}

// SetTimeout of the http requests
func (hc *hubHTTPClient) SetTimeout(t time.Duration) {
	hc.httpClient.Timeout = t
}

// handleResponse reads http response body into byte slice
// if response contains an unexpected status code, error is returned
func handleResponse(resp *http.Response, inErr error) (b []byte, response *http.Response, err error) {
	if inErr != nil {
		return nil, nil, inErr
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = cerr
		}
	}()

	response = resp
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if !isOKResponseCode(resp.StatusCode) {
		return nil, response, fmt.Errorf("got unexpected response status code: %d. response: %s", resp.StatusCode, string(b))
	}

	if len(b) == 0 {
		return []byte(fmt.Sprintf("Response status: %s", resp.Status)), response, nil
	}

	return
}

// isOKResponseCode identifies whether provided
// response code matches the expected OK code
func isOKResponseCode(code int) bool {
	return code >= http.StatusOK && code < http.StatusMultipleChoices
}
