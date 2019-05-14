package notificationhubs

import (
	"time"
)

type (
	// Headers structure
	Headers map[string]string

	// Registration is a device registration to the hub
	Registration struct {
		DeviceID           string             `json:"deviceID,omitempty"`
		ExpirationTime     *time.Time         `json:"expirationTime,omitempty"`
		NotificationFormat NotificationFormat `json:"service,omitempty"`
		RegistrationID     string             `json:"registrationID,omitempty"`
		Tags               string             `json:"tags,omitempty"`
	}

	// TemplateRegistration is a device registration to the hub supporting a template
	TemplateRegistration struct {
		DeviceID       string         `json:"deviceID,omitempty"`
		ExpirationTime *time.Time     `json:"expirationTime,omitempty"`
		RegistrationID string         `json:"registrationID,omitempty"`
		Tags           string         `json:"tags,omitempty"`
		Platform       TargetPlatform `json:"platform,omitempty"`
		Template       string         `json:"template,omitempty"`
	}

	// Registrations is a list of RegistrationResults
	Registrations struct {
		ID      string               `xml:"id"      json:"id,omitempty"`
		Title   string               `xml:"title"   json:"title,omitempty"`
		Updated *time.Time           `xml:"updated" json:"updated,omitempty"`
		Entries []RegistrationResult `xml:"entry"   json:"entries,omitempty"`
	}

	// RegistrationResult is the response from registration
	RegistrationResult struct {
		ID                  string               `xml:"id"        json:"id,omitempty"`
		Published           *time.Time           `xml:"published" json:"published,omitempty"`
		RegistrationContent *RegistrationContent `xml:"content"   json:"content,omitempty"`
		Title               string               `xml:"title"     json:"title,omitempty"`
		Updated             *time.Time           `xml:"updated"   json:"updated,omitempty"`
	}

	// RegistrationContent is information about a specific device registration
	RegistrationContent struct {
		Format           NotificationFormat `xml:"-" json:"format,omitempty"`
		Target           TargetPlatform     `xml:"-" json:"target,omitempty"`
		RegisteredDevice *RegisteredDevice  `xml:"-" json:"registeredDevice,omitempty"`

		AppleRegistrationDescription         *RegisteredDevice `xml:"AppleRegistrationDescription"         json:"-"`
		AppleTemplateRegistrationDescription *RegisteredDevice `xml:"AppleTemplateRegistrationDescription" json:"-"`
		GcmRegistrationDescription           *RegisteredDevice `xml:"GcmRegistrationDescription"           json:"-"`
		GcmTemplateRegistrationDescription   *RegisteredDevice `xml:"GcmTemplateRegistrationDescription"   json:"-"`
	}

	// RegisteredDevice is a device registration to the hub
	RegisteredDevice struct {
		DeviceID       string     `xml:"-"              json:"deviceID,omitempty"`
		ETag           string     `xml:"ETag"           json:"eTag,omitempty"`
		ExpirationTime *time.Time `xml:"-"              json:"expirationTime,omitempty"`
		Template       string     `xml:"BodyTemplate"   json:"template,omitempty"`
		RegistrationID string     `xml:"RegistrationId" json:"registrationID,omitempty"`
		Tags           []string   `xml:"-"              json:"tags,omitempty"`

		DeviceToken          *string `xml:"DeviceToken"       json:"-"`
		ExpirationTimeString *string `xml:"ExpirationTime"    json:"-"`
		GcmRegistrationID    *string `xml:"GcmRegistrationId" json:"-"`
		TagsString           *string `xml:"Tags"              json:"-"`
	}

	// NotificationDetails is the detailed information about a sent or scheduled message
	NotificationDetails struct {
		ID                string                `xml:"NotificationId"`
		State             NotificationState     `xml:"State"`
		EnqueueTime       string                `xml:"EnqueueTime"`
		StartTime         string                `xml:"StartTime"`
		EndTime           string                `xml:"EndTime"`
		Body              string                `xml:"NotificationBody"`
		TargetPlatforms   string                `xml:"TargetPlatforms"`
		ApnsOutcomeCounts *NotificationOutcomes `xml:"ApnsOutcomeCounts"`
		GcmOutcomeCounts  *NotificationOutcomes `xml:"GcmOutcomeCounts"`
	}

	// NotificationTelemetry is the id of a sent or scheduled message
	NotificationTelemetry struct {
		NotificationMessageID string `json:"notificationMessageID,omitempty"`
	}

	// NotificationOutcomes array of outcomes
	NotificationOutcomes struct {
		Outcomes []NotificationOutcome `xml:"Outcome"`
	}

	// NotificationOutcome name value pair for statistics
	NotificationOutcome struct {
		Name  NotificationOutcomeName `xml:"Name"`
		Count int                     `xml:"Count"`
	}

	// NotificationState is the state of the notification
	NotificationState string

	// NotificationFormat is the format of a notification
	NotificationFormat string

	// NotificationOutcomeName is a possible outcome of a notification
	NotificationOutcomeName string

	// TargetPlatform is the specific platform
	TargetPlatform string
)
