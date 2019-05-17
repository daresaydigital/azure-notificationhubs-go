package notificationhubs

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"path"
	"strings"
	"time"
)

// newRegistration initializes and returns a Notification pointer
func newRegistration(deviceID string, expirationTime *time.Time, notificationFormat NotificationFormat,
	registrationID string, tags string) *Registration {
	return &Registration{
		deviceID,
		expirationTime,
		notificationFormat,
		registrationID,
		tags,
	}
}

// newTemplateRegistration initializes and returns a TemplateNotification pointer
func newTemplateRegistration(deviceID string, expirationTime *time.Time, registrationID string, tags string,
	platform TargetPlatform, template string) *TemplateRegistration {
	return &TemplateRegistration{
		deviceID,
		expirationTime,
		registrationID,
		tags,
		platform,
		template,
	}
}

// Normalize normalizes all devices in the feed
func (r *Registrations) normalize() {
	for _, entry := range r.Entries {
		if entry.RegistrationContent != nil {
			entry.RegistrationContent.normalize()
		}
	}
}

// Normalize the registration result
func (r RegistrationResult) normalize() {
	if r.RegistrationContent != nil {
		r.RegistrationContent.normalize()
	}
}

// Normalize normalizes the different devices
func (r *RegistrationContent) normalize() {
	if r.AppleRegistrationDescription != nil || r.AppleTemplateRegistrationDescription != nil {
		if r.AppleTemplateRegistrationDescription != nil {
			r.Format = Template
			r.Target = AppleTemplatePlatform
			r.RegisteredDevice = r.AppleTemplateRegistrationDescription
		} else {
			r.Format = AppleFormat
			r.Target = ApplePlatform
			r.RegisteredDevice = r.AppleRegistrationDescription
		}
		r.RegisteredDevice.DeviceID = *r.RegisteredDevice.DeviceToken
		r.RegisteredDevice.DeviceToken = nil
		r.AppleRegistrationDescription = nil
		r.AppleTemplateRegistrationDescription = nil
	} else if r.GcmRegistrationDescription != nil || r.GcmTemplateRegistrationDescription != nil {
		if r.GcmTemplateRegistrationDescription != nil {
			r.Format = Template
			r.Target = GcmTemplatePlatform
			r.RegisteredDevice = r.GcmTemplateRegistrationDescription
		} else {
			r.Format = GcmFormat
			r.Target = GcmPlatform
			r.RegisteredDevice = r.GcmRegistrationDescription
		}
		r.RegisteredDevice.DeviceID = *r.RegisteredDevice.GcmRegistrationID
		r.RegisteredDevice.GcmRegistrationID = nil
		r.GcmRegistrationDescription = nil
		r.GcmTemplateRegistrationDescription = nil
	}
	if r.RegisteredDevice != nil {
		expirationTime, err := time.Parse("2006-01-02T15:04:05.000Z", *r.RegisteredDevice.ExpirationTimeString)
		if err != nil { // The API just forwards the date string used by Apple, Google etc unfortunately. So format varies.
			expirationTime, _ = time.Parse("2006-01-02T15:04:05.000", *r.RegisteredDevice.ExpirationTimeString)
		}
		r.RegisteredDevice.ExpirationTime = &expirationTime
		r.RegisteredDevice.ExpirationTimeString = nil
		r.RegisteredDevice.Tags = strings.Split(*r.RegisteredDevice.TagsString, ",")
		r.RegisteredDevice.TagsString = nil
	}
}

// Registration reads one specific registration
func (h *NotificationHub) Registration(ctx context.Context, deviceID string) (raw []byte, registrationResult *RegistrationResult, err error) {
	var (
		regURL = h.generateAPIURL("registrations")
	)
	regURL.Path = path.Join(regURL.Path, deviceID)
	raw, _, err = h.exec(ctx, getMethod, regURL, Headers{}, nil)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(raw, &registrationResult); err != nil {
		return
	}
	registrationResult.RegistrationContent.normalize()
	return
}

// Registrations reads all registrations
func (h *NotificationHub) Registrations(ctx context.Context) (raw []byte, registrations *Registrations, err error) {
	raw, _, err = h.exec(ctx, getMethod, h.generateAPIURL("registrations"), Headers{}, nil)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(raw, &registrations); err != nil {
		return
	}
	registrations.normalize()
	return
}

// Register sends a device registration to the Azure hub
func (h *NotificationHub) Register(ctx context.Context, r Registration) (raw []byte, registrationResult *RegistrationResult, err error) {
	var (
		regURL  = h.generateAPIURL("registrations")
		method  = postMethod
		payload = ""
		headers = map[string]string{
			"Content-Type": "application/atom+xml;type=entry;charset=utf-8",
		}
	)

	switch r.NotificationFormat {
	case AppleFormat:
		payload = strings.Replace(appleRegXMLString, "{{DeviceID}}", r.DeviceID, 1)
	case GcmFormat:
		payload = strings.Replace(gcmRegXMLString, "{{DeviceID}}", r.DeviceID, 1)
	default:
		return nil, nil, errors.New("Notification format not implemented")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)

	if r.RegistrationID != "" {
		method = putMethod
		regURL.Path = path.Join(regURL.Path, r.RegistrationID)
	}

	raw, _, err = h.exec(ctx, method, regURL, headers, bytes.NewBufferString(payload))

	if err == nil {
		if err = xml.Unmarshal(raw, &registrationResult); err != nil {
			return
		}
	}
	registrationResult.normalize()
	return
}

// RegisterWithTemplate sends a device registration with template to the Azure hub
func (h *NotificationHub) RegisterWithTemplate(ctx context.Context, r TemplateRegistration) (raw []byte, registrationResult *RegistrationResult, err error) {
	var (
		regURL  = h.generateAPIURL("registrations")
		method  = postMethod
		payload = ""
		headers = map[string]string{
			"Content-Type": "application/atom+xml;type=entry;charset=utf-8",
		}
	)

	switch r.Platform {
	case ApplePlatform:
		payload = strings.Replace(appleTemplateRegXMLString, "{{DeviceID}}", r.DeviceID, 1)
	case GcmPlatform:
		payload = strings.Replace(gcmTemplateRegXMLString, "{{DeviceID}}", r.DeviceID, 1)
	default:
		return nil, nil, errors.New("Notification format not implemented")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)
	payload = strings.Replace(payload, "{{Template}}", r.Template, 1)

	if r.RegistrationID != "" {
		method = putMethod
		regURL.Path = path.Join(regURL.Path, r.RegistrationID)
	}

	raw, _, err = h.exec(ctx, method, regURL, headers, bytes.NewBufferString(payload))

	if err == nil {
		if err = xml.Unmarshal(raw, &registrationResult); err != nil {
			return
		}
	}
	if registrationResult != nil {
		registrationResult.normalize()
	}
	return
}
