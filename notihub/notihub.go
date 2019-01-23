/*
Package notihub represents an http client
for microsoft azure notification hub
*/
package notihub

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"gopkg.in/xmlpath.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	apiVersionParam = "api-version"
	apiVersionValue = "2015-01"
	scheme          = "https"
)

const (
	Template           NotificationFormat = "template"
	AndroidFormat      NotificationFormat = "gcm"
	AppleFormat        NotificationFormat = "apple"
	BaiduFormat        NotificationFormat = "baidu"
	KindleFormat       NotificationFormat = "adm"
	WindowsFormat      NotificationFormat = "windows"
	WindowsPhoneFormat NotificationFormat = "windowsphone"

	AppleRegTemplate string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
    <content type="application/xml">
        <AppleRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
            <Tags>{{Tags}}</Tags>
            <DeviceToken>{{DeviceId}}</DeviceToken>
        </AppleRegistrationDescription>
    </content>
</entry>`
	AndroidRegTemplate string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
    <content type="application/xml">
        <GcmRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
            <Tags>{{Tags}}</Tags>
            <GcmRegistrationId>{{DeviceId}}</GcmRegistrationId>
        </GcmRegistrationDescription>
    </content>
</entry>`
)

type (
	Notification struct {
		Format  NotificationFormat
		Payload []byte
	}

	NotificationFormat string

	Registration struct {
		RegistrationId string             `json:"registrationId"`
		DeviceId       string             `json:"deviceId"`
		Service        NotificationFormat `json:"service"`
		Tags           string             `json:"tags"`
		ExpirationTime *time.Time         `json:"expirationTime,omitmepty"`
	}

	RegistrationRes struct {
		Id             string    `xml:"id"`
		Title          string    `xml:"title"`
		Updated        time.Time `xml:"updated"`
		RegistrationId string
		ETag           string
		ExpirationTime time.Time
	}

	NotificationHub struct {
		sasKeyValue             string
		sasKeyName              string
		host                    string
		regPath                 string
		stdURL                  *url.URL
		scheduleURL             *url.URL
		client                  HubClient
		expirationTimeGenerator expirationTimeGenerator

		regIdPath *xmlpath.Path
		eTagPath  *xmlpath.Path
		expTmPath *xmlpath.Path
	}

	HubClient interface {
		Exec(req *http.Request) ([]byte, error)
	}

	expirationTimeGenerator interface {
		GenerateTimestamp() int64
	}

	expirationTimeGeneratorFunc func() int64

	hubHttpClient struct {
		httpClient *http.Client
	}
)

// GenerateTimestamp calls f()
func (f expirationTimeGeneratorFunc) GenerateTimestamp() int64 {
	return f()
}

// Exec executes notification hub http request and handles the response
func (hc *hubHttpClient) Exec(req *http.Request) ([]byte, error) {
	return handleResponse(hc.httpClient.Do(req))
}

// GetContentType returns Content-Type
// associated with NotificationFormat
func (f NotificationFormat) GetContentType() string {
	switch f {
	case Template,
		AppleFormat,
		AndroidFormat,
		KindleFormat,
		BaiduFormat:
		return "application/json"
	}

	return "application/xml"
}

// IsValid identifies whether notification format is valid
func (f NotificationFormat) IsValid() bool {
	return f == Template ||
		f == AndroidFormat ||
		f == AppleFormat ||
		f == BaiduFormat ||
		f == KindleFormat ||
		f == WindowsFormat ||
		f == WindowsPhoneFormat
}

// NewNotification initalizes and returns Notification pointer
func NewNotification(format NotificationFormat, payload []byte) (*Notification, error) {
	if !format.IsValid() {
		return nil, fmt.Errorf("unknown format '%s'", format)
	}

	return &Notification{format, payload}, nil
}

// String returns Notification string representation
func (n *Notification) String() string {
	return fmt.Sprintf("&{%s %s}", n.Format, string(n.Payload))
}

// NewNotificationHub initializes and returns NotificationHub pointer
func NewNotificationHub(connectionString, hubPath string) *NotificationHub {
	connData := strings.Split(connectionString, ";")

	hub := &NotificationHub{}
	var endpoint string

	for _, connItem := range connData {
		if len(connItem) >= 14 && connItem[:8] == "Endpoint" {
			endpoint = strings.Trim(connItem[14:], "/")
			continue
		}

		if len(connItem) >= 20 && connItem[:19] == "SharedAccessKeyName" {
			hub.sasKeyName = connItem[20:]
			continue
		}

		if len(connItem) >= 16 && connItem[:15] == "SharedAccessKey" {
			hub.sasKeyValue = connItem[16:]
			continue
		}
	}

	query := url.Values{apiVersionParam: {apiVersionValue}}

	hub.host = endpoint
	hub.stdURL = &url.URL{
		Host:     endpoint,
		Scheme:   scheme,
		Path:     path.Join(hubPath, "messages"),
		RawQuery: query.Encode(),
	}
	hub.scheduleURL = &url.URL{
		Host:     endpoint,
		Scheme:   scheme,
		Path:     path.Join(hubPath, "schedulednotifications"),
		RawQuery: query.Encode(),
	}

	hub.regPath = path.Join(hubPath, "registrations")

	hub.client = &hubHttpClient{&http.Client{}}
	hub.expirationTimeGenerator = expirationTimeGeneratorFunc(generateExpirationTimestamp)

	hub.regIdPath = xmlpath.MustCompile("/entry/content/*/RegistrationId")
	hub.eTagPath = xmlpath.MustCompile("/entry/content/*/ETag")
	hub.expTmPath = xmlpath.MustCompile("/entry/content/*/ExpirationTime")

	return hub
}

// Send publishes notification to the azure hub
func (h *NotificationHub) Send(n *Notification, orTags []string) ([]byte, error) {
	b, err := h.send(n, orTags, nil)
	if err != nil {
		return nil, fmt.Errorf("NotificationHub.Send: %s", err)
	}

	return b, nil
}

// Schedule pusblishes a scheduled notification to azure notification hub
func (h *NotificationHub) Schedule(n *Notification, orTags []string, deliverTime time.Time) ([]byte, error) {
	b, err := h.send(n, orTags, &deliverTime)
	if err != nil {
		return nil, fmt.Errorf("NotificationHub.Schedule: %s", err)
	}

	return b, nil
}

// send sends notification to the azure hub
func (h *NotificationHub) send(n *Notification, orTags []string, deliverTime *time.Time) ([]byte, error) {
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

	urlStr := h.stdURL.String()
	if deliverTime != nil && deliverTime.Unix() > time.Now().Unix() {
		urlStr = h.scheduleURL.String()
		headers["ServiceBusNotification-ScheduleTime"] = deliverTime.Format("2006-01-02T15:04:05")
	}

	req, err := http.NewRequest("POST", urlStr, buf)
	if err != nil {
		return nil, err
	}

	for header, val := range headers {
		req.Header.Set(header, val)
	}

	return h.client.Exec(req)
}

// generateSasToken generates and returns
// azure notification hub shared access signatue token
func (h *NotificationHub) generateSasToken() string {
	uri := &url.URL{Scheme: scheme, Host: h.host}
	targetUri := strings.ToLower(uri.String())

	expires := h.expirationTimeGenerator.GenerateTimestamp()
	toSign := fmt.Sprintf("%s\n%d", url.QueryEscape(targetUri), expires)

	mac := hmac.New(sha256.New, []byte(h.sasKeyValue))
	mac.Write([]byte(toSign))
	macb := mac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(macb)

	tokenParams := url.Values{
		"sr":  {targetUri},
		"sig": {signature},
		"se":  {fmt.Sprintf("%d", expires)},
		"skn": {h.sasKeyName},
	}

	return fmt.Sprintf("SharedAccessSignature %s", tokenParams.Encode())
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

// Register sends registration to the azure hub
func (h *NotificationHub) Register(r Registration) (RegistrationRes, []byte, error) {
	regRes := RegistrationRes{}
	token := h.generateSasToken()

	headers := map[string]string{
		"Authorization": token,
		"Content-Type":  "application/atom+xml;type=entry;charset=utf-8",
	}

	regPath := h.regPath
	method := "POST"
	payload := ""

	switch r.Service {
	case AppleFormat:
		payload = strings.Replace(AppleRegTemplate, "{{DeviceId}}", r.DeviceId, 1)
	case AndroidFormat:
		payload = strings.Replace(AndroidRegTemplate, "{{DeviceId}}", r.DeviceId, 1)
	default:
		return regRes, nil, errors.New("not implemented.")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)

	if r.RegistrationId != "" {
		regPath = path.Join(regPath, r.RegistrationId)
		method = "PUT"
	}
	query := url.Values{apiVersionParam: {apiVersionValue}}

	regURL := url.URL{Host: h.host, Scheme: scheme, Path: regPath, RawQuery: query.Encode()}
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
			if regId, ok := h.regIdPath.String(root); ok {
				regRes.RegistrationId = regId
			} else {
				return regRes, res, errors.New("RegistrationId not fount")
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
