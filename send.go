package notificationhubs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"
)

// Send publishes notification directly
// Format tags according to https://docs.microsoft.com/en-us/azure/notification-hubs/notification-hubs-tags-segment-push-message
// ex. "(follows_RedSox || follows_Cardinals) && location_Boston"
// or nil if no tags should be used
func (h *NotificationHub) Send(ctx context.Context, n *Notification, tags *string) (raw []byte, telemetry *NotificationTelemetry, err error) {
	raw, telemetry, err = h.send(ctx, n, tags, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("notificationhubs.SendDirect: %s", err)
	}
	return
}

// SendDirect publishes notification to a specific device
func (h *NotificationHub) SendDirect(ctx context.Context, n *Notification, deviceHandle string) (raw []byte, telemetry *NotificationTelemetry, err error) {
	raw, telemetry, err = h.sendDirect(ctx, n, deviceHandle)
	if err != nil {
		return nil, nil, fmt.Errorf("notificationhubs.SendDirect: %s", err)
	}
	return
}

// Schedule pusblishes a scheduled notification
// Format tags according to https://docs.microsoft.com/en-us/azure/notification-hubs/notification-hubs-tags-segment-push-message
// or nil if no tags should be used
func (h *NotificationHub) Schedule(ctx context.Context, n *Notification, tags *string, deliverTime time.Time) (raw []byte, telemetry *NotificationTelemetry, err error) {
	raw, telemetry, err = h.send(ctx, n, tags, &deliverTime)
	if err != nil {
		return nil, nil, fmt.Errorf("notificationhubs.Schedule: %s", err)
	}
	return
}

// send sends notification to the azure hub
func (h *NotificationHub) send(ctx context.Context, n *Notification, tags *string, deliverTime *time.Time) (raw []byte, telemetry *NotificationTelemetry, err error) {
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
			return nil, nil, errors.New("You can not schedule a notification in the past")
		}
	} else {
		_url.Path = path.Join(_url.Path, "messages")
	}

	raw, response, err := h.exec(ctx, postMethod, _url, headers, bytes.NewBuffer(n.Payload))
	if err != nil {
		return
	}
	telemetry, err = NewNotificationTelemetryFromHTTPResponse(response)
	return
}

func (h *NotificationHub) sendDirect(ctx context.Context, n *Notification, deviceHandle string) (raw []byte, telemetry *NotificationTelemetry, err error) {
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
	raw, response, err := h.exec(ctx, postMethod, _url, headers, bytes.NewBuffer(n.Payload))
	if err != nil {
		return
	}
	telemetry, err = NewNotificationTelemetryFromHTTPResponse(response)
	return
}
