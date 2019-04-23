package notificationhubs

import (
	"strings"
	"time"
)

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
	// GcmFormat (ted) notification
	GcmFormat NotificationFormat = "gcm"
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

	// GcmRegTemplate is the XML string for registering an iOS device
	// Replace {{Tags}} and {{DeviceID}} with the correct values
	GcmRegTemplate string = `<?xml version="1.0" encoding="utf-8"?>
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
		RegistrationID     string             `json:"registrationId"`
		DeviceID           string             `json:"deviceId"`
		NotificationFormat NotificationFormat `json:"service"`
		Tags               string             `json:"tags"`
		ExpirationTime     time.Time          `json:"expirationTime,omitempty"`
	}

	// Registrations is a list of RegistrationResults
	Registrations struct {
		Title   string               `xml:"title" json:"title"`
		ID      string               `xml:"id" json:"id"`
		Updated time.Time            `xml:"updated" json:"updated"`
		Entries []RegistrationResult `xml:"entry" json:"entry"`
	}

	// RegistrationResult is the response from registration
	RegistrationResult struct {
		ID                  string              `xml:"id" json:"id"`
		Title               string              `xml:"title" json:"title"`
		Updated             time.Time           `xml:"updated" json:"updated"`
		Published           time.Time           `xml:"published" json:"published"`
		RegistrationContent RegistrationContent `xml:"content" json:"content"`
	}

	// RegistrationContent is information about a specific device registration
	RegistrationContent struct {
		GcmRegistrationDescription   *RegistratedDevice `xml:"GcmRegistrationDescription,omitempty" json:"-"`
		AppleRegistrationDescription *RegistratedDevice `xml:"AppleRegistrationDescription,omitempty" json:"-"`
		RegistratedDevice            *RegistratedDevice `xml:"-" json:"registratedDevice"`
		Format                       NotificationFormat `xml:"-" json:"format"`
	}

	// RegistratedDevice is a device registration to the hub
	RegistratedDevice struct {
		ExpirationTime    string   `xml:"ExpirationTime" json:"expirationTime,omitempty"`
		RegistrationID    string   `xml:"RegistrationId" json:"registrationID"`
		ETag              string   `xml:"ETag" json:"eTag"`
		DeviceToken       string   `xml:"DeviceToken" json:"-"`
		GcmRegistrationID string   `xml:"GcmRegistrationId" json:"-"`
		TagsString        string   `xml:"Tags" json:"-"`
		DeviceID          string   `xml:"-" json:"deviceId"`
		Tags              []string `xml:"-" json:"tags"`
	}
)

// Normalize normalizes the different devices
func (r *RegistrationContent) Normalize() {
	if r.AppleRegistrationDescription != nil {
		r.Format = AppleFormat
		r.RegistratedDevice = r.AppleRegistrationDescription
		r.AppleRegistrationDescription = nil
		r.RegistratedDevice.DeviceID = r.RegistratedDevice.DeviceToken
	} else if r.GcmRegistrationDescription != nil {
		r.Format = GcmFormat
		r.RegistratedDevice = r.GcmRegistrationDescription
		r.GcmRegistrationDescription = nil
		r.RegistratedDevice.DeviceID = r.RegistratedDevice.GcmRegistrationID
	}
	if r.RegistratedDevice != nil {
		r.RegistratedDevice.Tags = strings.Split(r.RegistratedDevice.TagsString, ",")
	}
}
