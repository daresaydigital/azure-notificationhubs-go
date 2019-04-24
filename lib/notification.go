package lib

import "fmt"

type (
	// Notification is a message that can be sent through the hub
	Notification struct {
		Format  NotificationFormat
		Payload []byte
	}
)

// NewNotification initializes and returns a Notification pointer
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
