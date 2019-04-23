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

	// Http methods
	postMethod = "POST"
	getMethod  = "GET"
	putMethod  = "PUT"
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
		DeviceID           string             `json:"deviceId,omitempty"`
		ExpirationTime     time.Time          `json:"expirationTime,omitempty"`
		NotificationFormat NotificationFormat `json:"service,omitempty"`
		RegistrationID     string             `json:"registrationID,omitempty"`
		Tags               string             `json:"tags,omitempty"`
	}

	// Registrations is a list of RegistrationResults
	Registrations struct {
		ID      string               `xml:"id"      json:"id,omitempty"`
		Title   string               `xml:"title"   json:"title,omitempty"`
		Updated time.Time            `xml:"updated" json:"updated,omitempty"`
		Entries []RegistrationResult `xml:"entry"   json:"entries,omitempty"`
	}

	// RegistrationResult is the response from registration
	RegistrationResult struct {
		ID                  string               `xml:"id"        json:"id,omitempty"`
		Published           time.Time            `xml:"published" json:"published,omitempty"`
		RegistrationContent *RegistrationContent `xml:"content"   json:"content,omitempty"`
		Title               string               `xml:"title"     json:"title,omitempty"`
		Updated             time.Time            `xml:"updated"   json:"updated,omitempty"`
	}

	// RegistrationContent is information about a specific device registration
	RegistrationContent struct {
		AppleRegistrationDescription *RegistratedDevice `xml:"AppleRegistrationDescription" json:"-"`
		Format                       NotificationFormat `xml:"-"                            json:"format,omitempty"`
		GcmRegistrationDescription   *RegistratedDevice `xml:"GcmRegistrationDescription"   json:"-"`
		RegistratedDevice            *RegistratedDevice `xml:"-"                            json:"registratedDevice,omitempty"`
	}

	// RegistratedDevice is a device registration to the hub
	RegistratedDevice struct {
		DeviceID             string    `xml:"-"                 json:"deviceID,omitempty"`
		DeviceToken          *string   `xml:"DeviceToken"       json:"-"`
		ETag                 string    `xml:"ETag"              json:"eTag,omitempty"`
		ExpirationTimeString *string   `xml:"ExpirationTime"    json:"-"`
		ExpirationTime       time.Time `xml:"-"                 json:"expirationTime,omitempty"`
		GcmRegistrationID    *string   `xml:"GcmRegistrationId" json:"-"`
		RegistrationID       string    `xml:"RegistrationId"    json:"registrationID,omitempty"`
		Tags                 []string  `xml:"-"                 json:"tags,omitempty"`
		TagsString           *string   `xml:"Tags"              json:"-"`
	}
)

// Normalize normalizes all devices in the feed
func (r *Registrations) normalize() {
	for _, entry := range r.Entries {
		if entry.RegistrationContent != nil {
			entry.RegistrationContent.normalize()
		}
	}
}

// Normalize normalizes the different devices
func (r *RegistrationContent) normalize() {
	if r.AppleRegistrationDescription != nil {
		r.Format = AppleFormat
		r.RegistratedDevice = r.AppleRegistrationDescription
		r.AppleRegistrationDescription = nil
		r.RegistratedDevice.DeviceID = *r.RegistratedDevice.DeviceToken
		r.RegistratedDevice.DeviceToken = nil
	} else if r.GcmRegistrationDescription != nil {
		r.Format = GcmFormat
		r.RegistratedDevice = r.GcmRegistrationDescription
		r.GcmRegistrationDescription = nil
		r.RegistratedDevice.DeviceID = *r.RegistratedDevice.GcmRegistrationID
		r.RegistratedDevice.GcmRegistrationID = nil
	}
	expirationTime, err := time.Parse("2006-01-02T15:04:05.000Z", *r.RegistratedDevice.ExpirationTimeString)
	if err != nil { // The API uses more than one date format unfortunately :(
		expirationTime, _ = time.Parse("2006-01-02T15:04:05.000", *r.RegistratedDevice.ExpirationTimeString)
	}
	r.RegistratedDevice.ExpirationTime = expirationTime
	r.RegistratedDevice.ExpirationTimeString = nil
	if r.RegistratedDevice != nil {
		r.RegistratedDevice.Tags = strings.Split(*r.RegistratedDevice.TagsString, ",")
		r.RegistratedDevice.TagsString = nil
	}
}
