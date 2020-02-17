package notificationhubs

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/daresaydigital/azure-notificationhubs-go/utils"
)

// NotificationHub is a client for sending messages through Azure Notification Hubs
type NotificationHub struct {
	SasKeyValue string
	SasKeyName  string
	HubURL      *url.URL

	client                  utils.HTTPClient
	expirationTimeGenerator utils.ExpirationTimeGenerator
}

// newNotificationHub initializes and returns NotificationHub pointer
func newNotificationHub(connectionString, hubPath string) *NotificationHub {
	var (
		connData    = strings.Split(connectionString, ";")
		_url        = &url.URL{}
		sasKeyName  = ""
		sasKeyValue = ""
	)
	for _, connItem := range connData {
		if strings.HasPrefix(connItem, paramEndpoint) {
			hubURL, err := url.Parse(connItem[len(paramEndpoint):])
			if err == nil {
				_url = hubURL
			}
			continue
		}

		if strings.HasPrefix(connItem, paramSaasKeyName) {
			sasKeyName = connItem[len(paramSaasKeyName):]
			continue
		}

		if strings.HasPrefix(connItem, paramSaasKeyValue) {
			sasKeyValue = connItem[len(paramSaasKeyValue):]
			continue
		}
	}

	if _url.Scheme == schemeServiceBus || _url.Scheme == "" {
		_url.Scheme = schemeDefault
	}

	_url.Path = hubPath
	_url.RawQuery = url.Values{apiVersionParam: {apiVersionValue}}.Encode()
	return &NotificationHub{
		SasKeyName:  sasKeyName,
		SasKeyValue: sasKeyValue,
		HubURL:      _url,

		client:                  utils.NewHubHTTPClient(),
		expirationTimeGenerator: utils.NewExpirationTimeGenerator(),
	}
}

// SetHTTPClient makes it possible to use a custom http client
func (h *NotificationHub) SetHTTPClient(c utils.HTTPClient) {
	h.client = c
}

// OnRequest adds an optional hook to add more logging or other upon a request from the hub
func (h *NotificationHub) OnRequest(fun *utils.OnRequestFunc) {
	h.client.OnRequest(fun)
}

// SetExpirationTimeGenerator makes is possible to use a custom generator
func (h *NotificationHub) SetExpirationTimeGenerator(e utils.ExpirationTimeGenerator) {
	h.expirationTimeGenerator = e
}

// generateSasToken generates and returns
// azure notification hub shared access signature token
func (h *NotificationHub) generateSasToken() string {
	uri := &url.URL{
		Host:   h.HubURL.Host,
		Scheme: h.HubURL.Scheme,
	}
	targetURI := strings.ToLower(uri.String())

	expires := h.expirationTimeGenerator.GenerateTimestamp()
	toSign := fmt.Sprintf("%s\n%d", url.QueryEscape(targetURI), expires)

	mac := hmac.New(sha256.New, []byte(h.SasKeyValue))
	mac.Write([]byte(toSign))
	macb := mac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(macb)

	tokenParams := url.Values{
		"sr":  {targetURI},
		"sig": {signature},
		"se":  {fmt.Sprintf("%d", expires)},
		"skn": {h.SasKeyName},
	}

	return fmt.Sprintf("SharedAccessSignature %s", tokenParams.Encode())
}

// exec request using method to url
func (h *NotificationHub) exec(ctx context.Context, method string, url *url.URL, headers Headers, buf io.Reader) ([]byte, *http.Response, error) {
	headers["Authorization"] = h.generateSasToken()
	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)
	for header, val := range headers {
		req.Header.Set(header, val)
	}
	return h.client.Exec(req)
}

// generate an URL for path
func (h *NotificationHub) generateAPIURL(endpoint string) *url.URL {
	return &url.URL{
		Host:     h.HubURL.Host,
		Scheme:   h.HubURL.Scheme,
		Path:     path.Join(h.HubURL.Path, endpoint),
		RawQuery: h.HubURL.RawQuery,
	}
}
