package notificationhubs

import "time"

const (
	Template           NotificationFormat = "template"
	AndroidFormat      NotificationFormat = "gcm"
	AppleFormat        NotificationFormat = "apple"
	BaiduFormat        NotificationFormat = "baidu"
	KindleFormat       NotificationFormat = "adm"
	WindowsFormat      NotificationFormat = "windows"
	WindowsPhoneFormat NotificationFormat = "windowsphone"

	AppleRegTemplate string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
  <content type="application/xml">
    <AppleRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
      <Tags>{{Tags}}</Tags>
      <DeviceToken>{{DeviceID}}</DeviceToken>
    </AppleRegistrationDescription>
  </content>
</entry>`
	AndroidRegTemplate string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
  <content type="application/xml">
    <GcmRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
      <Tags>{{Tags}}</Tags>
      <GcmRegistrationId>{{DeviceID}}</GcmRegistrationId>
    </GcmRegistrationDescription>
  </content>
</entry>`
)

type (
	// Registration is a device registration to the hub
	Registration struct {
		// RegistrationID id of the registration
		RegistrationID string `json:"registrationId"`
		// DeviceID id of the device
		DeviceID string `json:"deviceId"`
		// NotificationFormat of which type of notification that should be received
		NotificationFormat NotificationFormat `json:"service"`
		// Tags for the device
		Tags string `json:"tags"`
		// ExpirationTime
		ExpirationTime *time.Time `json:"expirationTime,omitmepty"`
	}

	// RegistrationResponse is the response from registration
	RegistrationResult struct {
		// ID
		ID string `xml:"id"`
		// Title
		Title string `xml:"title"`
		// Updated
		Updated time.Time `xml:"updated"`
		// RegistrationID
		RegistrationID string
		// ETag
		ETag string
		// ExpirationTime
		ExpirationTime time.Time
	}
)
