package notificationhubs

import (
	"github.com/daresaydigital/azure-notificationhubs-go/lib"
)

// NewNotificationHub initializes and returns NotificationHub pointer
func NewNotificationHub(connectionString, hubPath string) *lib.NotificationHub {
	return lib.NewNotificationHub(connectionString, hubPath)
}

// NewNotification initalizes and returns Notification pointer
func NewNotification(format lib.NotificationFormat, payload []byte) *lib.Notification {
	notification, err := lib.NewNotification(format, payload)
	if err != nil {
		return nil
	}
	return notification
}
