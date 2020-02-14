package notificationhubs

import (
	"encoding/json"
	"fmt"
)

type (
	// Notification is a message that can be sent through the hub
	Notification struct {
		Format  NotificationFormat
		Payload []byte
	}

	// IosBackgroundNotificationPayload is the payload required for a background notification
	IosBackgroundNotificationPayload struct {
		Aps struct {
			ContentAvailable int `json:"content-available"`
		} `json:"aps"`
	}
)

// newNotification initializes and returns a Notification pointer
func newNotification(format NotificationFormat, payload []byte) (*Notification, error) {
	if !format.IsValid() {
		return nil, fmt.Errorf("unknown format '%s'", format)
	}

	return &Notification{format, payload}, nil
}

// String returns Notification string representation
func (n *Notification) String() string {
	return fmt.Sprintf("&{%s %s}", n.Format, string(n.Payload))
}

func isIosBackgroundNotification(payload []byte) bool {
	var backgroundNotification IosBackgroundNotificationPayload
	err := json.Unmarshal(payload, &backgroundNotification)
	if err != nil {
		return false
	}

	return backgroundNotification.Aps.ContentAvailable == 1
}
