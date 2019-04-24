package lib

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/daresaydigital/azure-notificationhubs-go/lib/utils"
)

// NotificationHub is a client for sending messages through Azure Notification Hubs
type NotificationHub struct {
	SasKeyValue string
	SasKeyName  string
	HubURL      *url.URL

	client                  utils.HTTPClient
	expirationTimeGenerator utils.ExpirationTimeGenerator
}

// NewNotificationHub initializes and retubrns NotificationHub pointer
func NewNotificationHub(connectionString, hubPath string) *NotificationHub {
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

// SetExpirationTimeGenerator makes is possible to use a custom generator
func (h *NotificationHub) SetExpirationTimeGenerator(e utils.ExpirationTimeGenerator) {
	h.expirationTimeGenerator = e
}

// Send publishes notification directly
// Format tags according to https://docs.microsoft.com/en-us/azure/notification-hubs/notification-hubs-tags-segment-push-message
// ex. "(follows_RedSox || follows_Cardinals) && location_Boston"
// or nil if no tags should be used
func (h *NotificationHub) Send(ctx context.Context, n *Notification, tags *string) ([]byte, error) {
	b, err := h.send(ctx, n, tags, nil)
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
// Format tags according to https://docs.microsoft.com/en-us/azure/notification-hubs/notification-hubs-tags-segment-push-message
// or nil if no tags should be used
func (h *NotificationHub) Schedule(ctx context.Context, n *Notification, tags *string, deliverTime time.Time) ([]byte, error) {
	b, err := h.send(ctx, n, tags, &deliverTime)
	if err != nil {
		return nil, fmt.Errorf("notificationHub.Schedule: %s", err)
	}
	return b, nil
}

// Registration reads one specific registration
func (h *NotificationHub) Registration(ctx context.Context, deviceID string) (*RegistrationResult, []byte, error) {
	var (
		res    = &RegistrationResult{}
		regURL = h.generateAPIURL("registrations")
	)
	regURL.Path = path.Join(regURL.Path, deviceID)
	rawResponse, err := h.exec(ctx, getMethod, regURL, Headers{}, nil)
	if err != nil {
		return nil, rawResponse, err
	}
	if err = xml.Unmarshal(rawResponse, &res); err != nil {
		return nil, rawResponse, err
	}
	res.RegistrationContent.normalize()
	return res, rawResponse, nil
}

// Registrations reads all registrations
func (h *NotificationHub) Registrations(ctx context.Context) (*Registrations, []byte, error) {
	rawResponse, err := h.exec(ctx, getMethod, h.generateAPIURL("registrations"), Headers{}, nil)
	if err != nil {
		return nil, rawResponse, err
	}
	res := &Registrations{}
	if err = xml.Unmarshal(rawResponse, &res); err != nil {
		return nil, rawResponse, err
	}
	res.normalize()
	return res, rawResponse, nil
}

// Register sends registration to the azure hub
func (h *NotificationHub) Register(ctx context.Context, r Registration) (RegistrationResult, []byte, error) {
	var (
		regRes  = RegistrationResult{}
		regURL  = h.generateAPIURL("registrations")
		method  = postMethod
		payload = ""
		headers = map[string]string{
			"Content-Type": "application/atom+xml;type=entry;charset=utf-8",
		}
	)

	switch r.NotificationFormat {
	case AppleFormat:
		payload = strings.Replace(AppleRegTemplate, "{{DeviceID}}", r.DeviceID, 1)
	case GcmFormat:
		payload = strings.Replace(GcmRegTemplate, "{{DeviceID}}", r.DeviceID, 1)
	default:
		return regRes, nil, errors.New("Notification format not implemented")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)

	if r.RegistrationID != "" {
		method = putMethod
		regURL.Path = path.Join(regURL.Path, r.RegistrationID)
	}

	res, err := h.exec(ctx, method, regURL, headers, bytes.NewBufferString(payload))

	if err == nil {
		if err = xml.Unmarshal(res, &regRes); err != nil {
			return regRes, res, err
		}
		regRes.RegistrationContent.normalize()
	}
	return regRes, res, err
}

// send sends notification to the azure hub
func (h *NotificationHub) send(ctx context.Context, n *Notification, tags *string, deliverTime *time.Time) ([]byte, error) {
	var (
		headers = map[string]string{
			"Content-Type":                  n.Format.GetContentType(),
			"ServiceBusNotification-Format": string(n.Format),
			"X-Apns-Expiration":             string(h.expirationTimeGenerator.GenerateTimestamp()), //apns-expiration
		}
		_url = h.generateAPIURL("")
	)

	if tags != nil && len(*tags) > 0 {
		headers["ServiceBusNotification-Tags"] = *tags
	}

	if deliverTime != nil {
		if deliverTime.After(time.Now()) {
			_url.Path = path.Join(_url.Path, "schedulednotifications")
			headers["ServiceBusNotification-ScheduleTime"] = deliverTime.Format("2006-01-02T15:04:05")
		} else {
			return nil, errors.New("You can not schedule a notification in the past")
		}
	} else {
		_url.Path = path.Join(_url.Path, "messages")
	}

	return h.exec(ctx, postMethod, _url, headers, bytes.NewBuffer(n.Payload))
}

func (h *NotificationHub) sendDirect(ctx context.Context, n *Notification, deviceHandle string) ([]byte, error) {
	var (
		headers = Headers{
			"Content-Type":                        n.Format.GetContentType(),
			"ServiceBusNotification-Format":       string(n.Format),
			"ServiceBusNotification-DeviceHandle": deviceHandle,
			"X-Apns-Expiration":                   string(h.expirationTimeGenerator.GenerateTimestamp()), //apns-expiration
		}
		query = h.HubURL.Query()
	)
	query.Add(directParam, "")
	_url := &url.URL{
		Host:     h.HubURL.Host,
		Scheme:   h.HubURL.Scheme,
		Path:     path.Join(h.HubURL.Path, "messages"),
		RawQuery: query.Encode(),
	}
	return h.exec(ctx, postMethod, _url, headers, bytes.NewBuffer(n.Payload))
}

// generateSasToken generates and returns
// azure notification hub shared access signatue token
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
func (h *NotificationHub) exec(ctx context.Context, method string, url *url.URL, headers Headers, buf io.Reader) ([]byte, error) {
	headers["Authorization"] = h.generateSasToken()
	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
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
