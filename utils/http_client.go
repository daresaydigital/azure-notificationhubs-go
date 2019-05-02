package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// HTTPClient interface, replaceable for testing or custom implementation
	HTTPClient interface {
		Exec(req *http.Request) ([]byte, *http.Response, error)
	}

	// HubHTTPClient is the internal HTTPClient
	HubHTTPClient struct {
		httpClient *http.Client
	}
)

// NewHubHTTPClient is creating the default client
func NewHubHTTPClient() HTTPClient {
	return HubHTTPClient{
		httpClient: &http.Client{},
	}
}

// Exec executes notification hub http request and handles the response
func (hc HubHTTPClient) Exec(req *http.Request) ([]byte, *http.Response, error) {
	return handleResponse(hc.httpClient.Do(req))
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
		return nil, response, fmt.Errorf("Got unexpected response status code: %d. response: %s", resp.StatusCode, string(b))
	}

	if len(b) == 0 {
		return []byte(fmt.Sprintf("Response status: %s", resp.Status)), response, nil
	}

	return
}

// isOKResponseCode identifies whether provided
// response code matches the expected OK code
func isOKResponseCode(code int) bool {
	return code == http.StatusCreated || code == http.StatusOK
}
