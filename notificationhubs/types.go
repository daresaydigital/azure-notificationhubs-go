package notificationhubs

import "time"

// Internal constants
const (
	apiVersionParam = "api-version"
	apiVersionValue = "2015-01" // Looks old but the API is the same
	directParam     = "direct"

	// for connection string parsing
	schemeServiceBus  = "sb"
	schemeDefault     = "https"
	paramEndpoint     = "Endpoint="
	paramSaasKeyName  = "SharedAccessKeyName="
	paramSaasKeyValue = "SharedAccessKey="
)

// Public constants
const (
	// Template notification
	Template NotificationFormat = "template"
	// AndroidFormat (ted) notification
	AndroidFormat NotificationFormat = "gcm"
	// AppleFormat (ted) notification
	AppleFormat NotificationFormat = "apple"
	// BaiduFormat (ted) notification
	BaiduFormat NotificationFormat = "baidu"
	// KindleFormat (ted) notification
	KindleFormat NotificationFormat = "adm"
	// WindowsFormat (ted) notification
	WindowsFormat NotificationFormat = "windows"
	// WindowsPhoneFormat (ted) notification
	WindowsPhoneFormat NotificationFormat = "windowsphone"

	// AppleRegTemplate is the XML string for registering an iOS device
	// Replace {{Tags}} and {{DeviceID}} with the correct values
	AppleRegTemplate string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
  <content type="application/xml">
    <AppleRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
      <Tags>{{Tags}}</Tags>
      <DeviceToken>{{DeviceID}}</DeviceToken>
    </AppleRegistrationDescription>
  </content>
</entry>`

	// AndroidRegTemplate is the XML string for registering an iOS device
	// Replace {{Tags}} and {{DeviceID}} with the correct values
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
	// Headers structure
	Headers map[string]string

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

	// RegistrationResult is the response from registration
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

	// Registrations is a list of RegistrationResults
	Registrations struct {
		Feed struct {
			Title   string               `xml:"title"`
			Entries []RegistrationResult `xml:"entry"`
		} `xml:"feed"`
	}
)
