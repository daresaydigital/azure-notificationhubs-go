package notificationhubs

// NewNotificationHub initializes and returns NotificationHub pointer
func NewNotificationHub(connectionString, hubPath string) *NotificationHub {
	return newNotificationHub(connectionString, hubPath)
}

// NewNotification initalizes and returns Notification pointer
func NewNotification(format NotificationFormat, payload []byte) (*Notification, error) {
	return newNotification(format, payload)
}
