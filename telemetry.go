package notificationhubs

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"path"
	"regexp"
)

// NotificationDetails reads one specific registration
func (h *NotificationHub) NotificationDetails(ctx context.Context, notificationID string) (details *NotificationDetails, raw []byte, err error) {
	var (
		// res  = &NotificationDetails{}
		_url = h.generateAPIURL("messages")
	)
	_url.Path = path.Join(_url.Path, notificationID)
	_url.RawQuery = url.Values{apiVersionParam: {telemetryAPIVersionValue}}.Encode()
	raw, _, err = h.exec(ctx, getMethod, _url, Headers{}, nil)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(raw, &details); err != nil {
		return
	}
	return
}

// NewNotificationTelemetryFromLocationURL create Telemetry from Location URL
func NewNotificationTelemetryFromLocationURL(url string) *NotificationTelemetry {
	var re = regexp.MustCompile(`/messages/(?P<id>.*)\?api-version=`)
	groupNames := re.SubexpNames()
	for _, match := range re.FindAllStringSubmatch(url, -1) {
		for groupIdx, group := range match {
			name := groupNames[groupIdx]
			if name == "id" {
				return &NotificationTelemetry{
					NotificationMessageID: group,
				}
			}
		}
	}
	return nil
}

// NewNotificationTelemetryFromHTTPResponse reads the Location header from URL
// Notification Telemetry is only available for Standard tier Notification Hubs.
func NewNotificationTelemetryFromHTTPResponse(response *http.Response) (*NotificationTelemetry, error) {
	if response == nil || response.Header == nil {
		return nil, errors.New("Could not parse telemetry from response")
	}
	location := response.Header.Get("location")
	if len(location) == 0 {
		return &NotificationTelemetry{}, nil
	}
	return NewNotificationTelemetryFromLocationURL(location), nil
}
