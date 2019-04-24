package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// HTTPClient interfacce
	HTTPClient interface {
		Exec(req *http.Request) ([]byte, error)
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
func (hc HubHTTPClient) Exec(req *http.Request) ([]byte, error) {
	return handleResponse(hc.httpClient.Do(req))
}

// handleResponse reads http response body into byte slice
// if response contains an unexpected status code, error is returned
func handleResponse(resp *http.Response, inErr error) (b []byte, err error) {
	if inErr != nil {
		return nil, inErr
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = cerr
		}
	}()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !isOKResponseCode(resp.StatusCode) {
		return nil, fmt.Errorf("got unexpected response status code: %d. response: %s", resp.StatusCode, b)
	}

	if len(b) == 0 {
		return []byte(fmt.Sprintf("response status: %s", resp.Status)), nil
	}

	return
}

// isOKResponseCode identifies whether provided
// response code matches the expected OK code
func isOKResponseCode(code int) bool {
	return code == http.StatusCreated || code == http.StatusOK
}
