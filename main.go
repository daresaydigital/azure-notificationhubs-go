package azurenotificationhubs

import (
	nh "github.com/daresaydigital/azure-notificationhubs-go/notificationhubs"
)

// NewNotificationHub initializes and returns NotificationHub pointer
func NewNotificationHub(connectionString, hubPath string) *nh.NotificationHub {
	return nh.NewNotificationHub(connectionString, hubPath)
}

// NewNotification initalizes and returns Notification pointer
func NewNotification(format nh.NotificationFormat, payload []byte) *nh.Notification {
	notification, err := nh.NewNotification(format, payload)
	if err != nil {
		return nil
	}
	return notification
}
