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
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/daresaydigital/azure-notificationhubs-go/notificationhubs/utils"
	"gopkg.in/xmlpath.v2"
)

// NotificationHub is a client for sending messages through Azure Notification Hubs
type NotificationHub struct {
	sasKeyValue             string
	sasKeyName              string
	hubURL                  *url.URL
	client                  utils.HTTPClient
	expirationTimeGenerator utils.ExpirationTimeGenerator

	regIDPath *xmlpath.Path
	eTagPath  *xmlpath.Path
	expTmPath *xmlpath.Path
}

// SetHTTPClient makes it possible to use a custom http client
func (h *NotificationHub) SetHTTPClient(c utils.HTTPClient) {
	h.client = c
}

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

// Registrations reads all registrations
func (h *NotificationHub) Registrations(ctx context.Context) (*Registrations, []byte, error) {
	rawResponse, err := h.exec(ctx, "GET", h.generateAPIURL("registrations"), Headers{}, nil)
	if err != nil {
		return nil, rawResponse, err
	}
	result := &Registrations{}
	if err = xml.Unmarshal(rawResponse, &result); err != nil {
		return result, rawResponse, err
	}
	return result, rawResponse, nil
}

// Register sends registration to the azure hub
func (h *NotificationHub) Register(ctx context.Context, r Registration) (RegistrationResult, []byte, error) {
	var (
		regRes  = RegistrationResult{}
		regURL  = h.generateAPIURL("registrations")
		method  = "POST"
		payload = ""
		headers = map[string]string{
			"Content-Type": "application/atom+xml;type=entry;charset=utf-8",
		}
	)

	switch r.NotificationFormat {
	case AppleFormat:
		payload = strings.Replace(AppleRegTemplate, "{{DeviceID}}", r.DeviceID, 1)
	case AndroidFormat:
		payload = strings.Replace(AndroidRegTemplate, "{{DeviceID}}", r.DeviceID, 1)
	default:
		return regRes, nil, errors.New("Notification format not implemented")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)

	if r.RegistrationID != "" {
		method = "PUT"
		regURL.Path = path.Join(regURL.Path, r.RegistrationID)
	}

	res, err := h.exec(ctx, method, regURL, headers, bytes.NewBufferString(payload))

	if err == nil {
		if err = xml.Unmarshal(res, &regRes); err != nil {
			return regRes, res, err
		}
		rb := bytes.NewReader(res)
		if root, err := xmlpath.Parse(rb); err == nil {
			if regID, ok := h.regIDPath.String(root); ok {
				regRes.RegistrationID = regID
			} else {
				return regRes, res, errors.New("RegistrationID not found")
			}
			if etag, ok := h.eTagPath.String(root); ok {
				regRes.ETag = etag
			} else {
				return regRes, res, errors.New("ETag not found")
			}
			if expTm, ok := h.expTmPath.String(root); ok {
				if regRes.ExpirationTime, err = time.Parse("2006-01-02T15:04:05.999", expTm); err != nil {
					return regRes, res, err
				}
			} else {
				return regRes, res, err
			}
		} else {
			return regRes, res, errors.New("ExpirationTime not found")
		}
	}
	return regRes, res, err
}

// send sends notification to the azure hub
func (h *NotificationHub) send(ctx context.Context, n *Notification, orTags []string, deliverTime *time.Time) ([]byte, error) {
	var (
		headers = map[string]string{
			"Content-Type":                  n.Format.GetContentType(),
			"ServiceBusNotification-Format": string(n.Format),
			"X-Apns-Expiration":             string(h.expirationTimeGenerator.GenerateTimestamp()), //apns-expiration
		}
		_url = h.generateAPIURL("")
	)

	if len(orTags) > 0 {
		headers["ServiceBusNotification-Tags"] = strings.Join(orTags, " || ")
	}

	if deliverTime != nil && deliverTime.Unix() > time.Now().Unix() {
		_url.Path = path.Join(_url.Path, "schedulednotifications")
		headers["ServiceBusNotification-ScheduleTime"] = deliverTime.Format("2006-01-02T15:04:05")
	} else {
		_url.Path = path.Join(_url.Path, "messages")
	}

	return h.exec(ctx, "POST", _url, headers, bytes.NewBuffer(n.Payload))
}

func (h *NotificationHub) sendDirect(ctx context.Context, n *Notification, deviceHandle string) ([]byte, error) {
	var (
		headers = Headers{
			"Content-Type":                        n.Format.GetContentType(),
			"ServiceBusNotification-Format":       string(n.Format),
			"ServiceBusNotification-DeviceHandle": deviceHandle,
			"X-Apns-Expiration":                   string(h.expirationTimeGenerator.GenerateTimestamp()), //apns-expiration
		}
		query = h.hubURL.Query()
	)
	query.Add(directParam, "")
	_url := &url.URL{
		Host:     h.hubURL.Host,
		Scheme:   h.hubURL.Scheme,
		Path:     path.Join(h.hubURL.Path, "messages"),
		RawQuery: query.Encode(),
	}
	return h.exec(ctx, "POST", _url, headers, bytes.NewBuffer(n.Payload))
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
		Host:     h.hubURL.Host,
		Scheme:   h.hubURL.Scheme,
		Path:     path.Join(h.hubURL.Path, endpoint),
		RawQuery: h.hubURL.RawQuery,
	}
}
