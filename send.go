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
