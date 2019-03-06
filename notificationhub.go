package notificationhubs

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"gopkg.in/xmlpath.v2"
)

const (
	apiVersionParam = "api-version"
	apiVersionValue = "2015-01"
	directParam     = "direct"

	// for connection string parsing
	schemeServiceBus  = "sb"
	schemeDefault     = "https"
	paramEndpoint     = "Endpoint="
	paramSaasKeyName  = "SharedAccessKeyName="
	paramSaasKeyValue = "SharedAccessKey="
)

type (
	// NotificationHub is a client for sending messages through Azure Notification Hubs
	NotificationHub struct {
		sasKeyValue             string
		sasKeyName              string
		hubURL                  *url.URL
		client                  hubClient
		expirationTimeGenerator expirationTimeGenerator

		regIDPath *xmlpath.Path
		eTagPath  *xmlpath.Path
		expTmPath *xmlpath.Path
	}

	hubClient interface {
		Exec(req *http.Request) ([]byte, error)
	}

	expirationTimeGenerator interface {
		GenerateTimestamp() int64
	}

	expirationTimeGeneratorFunc func() int64

	hubHTTPClient struct {
		httpClient *http.Client
	}
)

// Send publishes notification
func (h *NotificationHub) Send(ctx context.Context, n *Notification, orTags []string) ([]byte, error) {
	b, err := h.send(ctx, n, orTags, nil)
	if err != nil {
		return nil, fmt.Errorf("notificationHub.Send: %s", err)
	}

	return b, nil
}

// SendDirect publishes notification to a specific device
func (h *NotificationHub) SendDirect(ctx context.Context, n *Notification, deviceHandle string) ([]byte, error) {
	b, err := h.sendDirect(ctx, n, deviceHandle)
	if err != nil {
		return nil, fmt.Errorf("notificationHub.SendDirect: %s", err)
	}

	return b, nil
}

// Schedule pusblishes a scheduled notification
func (h *NotificationHub) Schedule(ctx context.Context, n *Notification, orTags []string, deliverTime time.Time) ([]byte, error) {
	b, err := h.send(ctx, n, orTags, &deliverTime)
	if err != nil {
		return nil, fmt.Errorf("notificationHub.Schedule: %s", err)
	}

	return b, nil
}

// Register sends registration to the azure hub
func (h *NotificationHub) Register(r Registration) (RegistrationResult, []byte, error) {
	regRes := RegistrationResult{}
	token := h.generateSasToken()

	headers := map[string]string{
		"Authorization": token,
		"Content-Type":  "application/atom+xml;type=entry;charset=utf-8",
	}

	payload := ""

	switch r.NotificationFormat {
	case AppleFormat:
		payload = strings.Replace(AppleRegTemplate, "{{DeviceID}}", r.DeviceID, 1)
	case AndroidFormat:
		payload = strings.Replace(AndroidRegTemplate, "{{DeviceID}}", r.DeviceID, 1)
	default:
		return regRes, nil, errors.New("Not implemented")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)

	method := "POST"
	regURL := url.URL{
		Host:     h.hubURL.Host,
		Scheme:   h.hubURL.Scheme,
		Path:     path.Join(h.hubURL.Path, "registrations"),
		RawQuery: h.hubURL.RawQuery,
	}

	if r.RegistrationID != "" {
		method = "PUT"
		regURL.Path = path.Join(regURL.Path, r.RegistrationID)
	}

	urlStr := regURL.String()
	buf := bytes.NewBufferString(payload)
	req, err := http.NewRequest(method, urlStr, buf)
	if err != nil {
		return regRes, nil, err
	}

	for header, val := range headers {
		req.Header.Set(header, val)
	}

	res, err := h.client.Exec(req)
	if err == nil {
		if err = xml.Unmarshal(res, &regRes); err != nil {
			return regRes, res, err
		}
		rb := bytes.NewReader(res)
		if root, err := xmlpath.Parse(rb); err == nil {
			if regID, ok := h.regIDPath.String(root); ok {
				regRes.RegistrationID = regID
			} else {
				return regRes, res, errors.New("RegistrationID not fount")
			}
			if etag, ok := h.eTagPath.String(root); ok {
				regRes.ETag = etag
			} else {
				return regRes, res, errors.New("ETag not fount")
			}
			if expTm, ok := h.expTmPath.String(root); ok {
				if regRes.ExpirationTime, err = time.Parse("2006-01-02T15:04:05.999", expTm); err != nil {
					return regRes, res, err
				}
			} else {
				return regRes, res, err
			}
		} else {
			return regRes, res, errors.New("ExpirationTime not fount")
		}
	}
	return regRes, res, err
}

// send sends notification to the azure hub
func (h *NotificationHub) send(ctx context.Context, n *Notification, orTags []string, deliverTime *time.Time) ([]byte, error) {
	token := h.generateSasToken()
	buf := bytes.NewBuffer(n.Payload)

	headers := map[string]string{
		"Authorization":                 token,
		"Content-Type":                  n.Format.GetContentType(),
		"ServiceBusNotification-Format": string(n.Format),
		"X-Apns-Expiration":             string(generateExpirationTimestamp()), //apns-expiration
	}

	if len(orTags) > 0 {
		headers["ServiceBusNotification-Tags"] = strings.Join(orTags, " || ")
	}

	_url := &url.URL{
		Host:     h.hubURL.Host,
		Scheme:   h.hubURL.Scheme,
		Path:     h.hubURL.Path,
		RawQuery: h.hubURL.RawQuery,
	}

	if deliverTime != nil && deliverTime.Unix() > time.Now().Unix() {
		_url.Path = path.Join(_url.Path, "schedulednotifications")
		headers["ServiceBusNotification-ScheduleTime"] = deliverTime.Format("2006-01-02T15:04:05")
	} else {
		_url.Path = path.Join(_url.Path, "messages")
	}

	req, err := http.NewRequest("POST", _url.String(), buf)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	for header, val := range headers {
		req.Header.Set(header, val)
	}

	return h.client.Exec(req)
}

func (h *NotificationHub) sendDirect(ctx context.Context, n *Notification, deviceHandle string) ([]byte, error) {
	token := h.generateSasToken()
	buf := bytes.NewBuffer(n.Payload)

	headers := map[string]string{
		"Authorization":                       token,
		"Content-Type":                        n.Format.GetContentType(),
		"ServiceBusNotification-Format":       string(n.Format),
		"ServiceBusNotification-DeviceHandle": deviceHandle,
		"X-Apns-Expiration":                   string(generateExpirationTimestamp()), //apns-expiration
	}

	query := h.hubURL.Query()
	query.Add(directParam, "")

	_url := &url.URL{
		Host:     h.hubURL.Host,
		Scheme:   h.hubURL.Scheme,
		Path:     path.Join(h.hubURL.Path, "messages"),
		RawQuery: query.Encode(),
	}

	req, err := http.NewRequest("POST", _url.String(), buf)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	for header, val := range headers {
		req.Header.Set(header, val)
	}

	return h.client.Exec(req)
}

// generateSasToken generates and returns
// azure notification hub shared access signatue token
func (h *NotificationHub) generateSasToken() string {
	uri := &url.URL{
		Host:   h.hubURL.Host,
		Scheme: h.hubURL.Scheme,
	}
	targetURI := strings.ToLower(uri.String())

	expires := h.expirationTimeGenerator.GenerateTimestamp()
	toSign := fmt.Sprintf("%s\n%d", url.QueryEscape(targetURI), expires)

	mac := hmac.New(sha256.New, []byte(h.sasKeyValue))
	mac.Write([]byte(toSign))
	macb := mac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(macb)

	tokenParams := url.Values{
		"sr":  {targetURI},
		"sig": {signature},
		"se":  {fmt.Sprintf("%d", expires)},
		"skn": {h.sasKeyName},
	}

	return fmt.Sprintf("SharedAccessSignature %s", tokenParams.Encode())
}

// Exec executes notification hub http request and handles the response
func (hc *hubHTTPClient) Exec(req *http.Request) ([]byte, error) {
	return handleResponse(hc.httpClient.Do(req))
}

// GenerateTimestamp calls f()
func (f expirationTimeGeneratorFunc) GenerateTimestamp() int64 {
	return f()
}

// generateExpirationTimestamp generates token expiration timestamp value
func generateExpirationTimestamp() int64 {
	return time.Now().Unix() + int64(3600)
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
