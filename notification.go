package notificationhubs

import "fmt"

type (
	// Notification is a message that can be sent through the hub
	Notification struct {
		Format  NotificationFormat
		Payload []byte
	}
)

// String returns Notification string representation
func (n *Notification) String() string {
	return fmt.Sprintf("&{%s %s}", n.Format, string(n.Payload))
}
