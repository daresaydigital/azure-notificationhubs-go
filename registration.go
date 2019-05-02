package notificationhubs

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
)

// Normalize normalizes all devices in the feed
func (r *Registrations) normalize() {
	for _, entry := range r.Entries {
		if entry.RegistrationContent != nil {
			entry.RegistrationContent.normalize()
		}
	}
}

// Normalize the registration result
func (r *RegistrationResult) normalize() {
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
			r.RegistratedDevice = r.AppleTemplateRegistrationDescription
		} else {
			r.Format = AppleFormat
			r.Target = ApplePlatform
			r.RegistratedDevice = r.AppleRegistrationDescription
		}
		r.RegistratedDevice.DeviceID = *r.RegistratedDevice.DeviceToken
		r.RegistratedDevice.DeviceToken = nil
		r.AppleRegistrationDescription = nil
		r.AppleTemplateRegistrationDescription = nil
	} else if r.GcmRegistrationDescription != nil || r.GcmTemplateRegistrationDescription != nil {
		if r.GcmTemplateRegistrationDescription != nil {
			r.Format = Template
			r.Target = GcmTemplatePlatform
			r.RegistratedDevice = r.GcmTemplateRegistrationDescription
		} else {
			r.Format = GcmFormat
			r.Target = GcmPlatform
			r.RegistratedDevice = r.GcmRegistrationDescription
		}
		r.RegistratedDevice.DeviceID = *r.RegistratedDevice.GcmRegistrationID
		r.RegistratedDevice.GcmRegistrationID = nil
		r.GcmRegistrationDescription = nil
		r.GcmTemplateRegistrationDescription = nil
	}
	if r.RegistratedDevice != nil {
		expirationTime, err := time.Parse("2006-01-02T15:04:05.000Z", *r.RegistratedDevice.ExpirationTimeString)
		if err != nil { // The API just forwards the date string used by Apple, Google etc unfortunately. So format varies.
			expirationTime, _ = time.Parse("2006-01-02T15:04:05.000", *r.RegistratedDevice.ExpirationTimeString)
		}
		r.RegistratedDevice.ExpirationTime = &expirationTime
		r.RegistratedDevice.ExpirationTimeString = nil
		r.RegistratedDevice.Tags = strings.Split(*r.RegistratedDevice.TagsString, ",")
		r.RegistratedDevice.TagsString = nil
	}
}

// Registration reads one specific registration
func (h *NotificationHub) Registration(ctx context.Context, deviceID string) (*RegistrationResult, []byte, error) {
	var (
		res    = &RegistrationResult{}
		regURL = h.generateAPIURL("registrations")
	)
	regURL.Path = path.Join(regURL.Path, deviceID)
	rawResponse, _, err := h.exec(ctx, getMethod, regURL, Headers{}, nil)
	if err != nil {
		return nil, rawResponse, err
	}
	if err = xml.Unmarshal(rawResponse, &res); err != nil {
		return nil, rawResponse, err
	}
	res.RegistrationContent.normalize()
	return res, rawResponse, nil
}

// Registrations reads all registrations
func (h *NotificationHub) Registrations(ctx context.Context) (registrations *Registrations, rawResponse []byte, err error) {
	rawResponse, _, err = h.exec(ctx, getMethod, h.generateAPIURL("registrations"), Headers{}, nil)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(rawResponse, &registrations); err != nil {
		return
	}
	registrations.normalize()
	return
}

// Register sends a device registration to the Azure hub
func (h *NotificationHub) Register(ctx context.Context, r Registration) (regRes *RegistrationResult, raw []byte, err error) {
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
		if err = xml.Unmarshal(raw, &regRes); err != nil {
			return
		}
	}
	regRes.normalize()
	return
}

// RegisterWithTemplate sends a device registration with template to the Azure hub
func (h *NotificationHub) RegisterWithTemplate(ctx context.Context, r TemplateRegistration) (regRes *RegistrationResult, raw []byte, err error) {
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
		return regRes, nil, errors.New("Notification format not implemented")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)
	payload = strings.Replace(payload, "{{Template}}", r.Template, 1)

	if r.RegistrationID != "" {
		method = putMethod
		regURL.Path = path.Join(regURL.Path, r.RegistrationID)
	}

	raw, _, err = h.exec(ctx, method, regURL, headers, bytes.NewBufferString(payload))

	fmt.Printf("Raw: %s\n", string(raw))

	if err == nil {
		if err = xml.Unmarshal(raw, &regRes); err != nil {
			return
		}
	}
	regRes.normalize()
	return
}
