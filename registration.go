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

// Registration reads one specific registration
func (h *NotificationHub) Registration(ctx context.Context, deviceID string) (*RegistrationResult, []byte, error) {
	var (
		res    = &RegistrationResult{}
		regURL = h.generateAPIURL("registrations")
	)
	regURL.Path = path.Join(regURL.Path, deviceID)
	rawResponse, err := h.exec(ctx, getMethod, regURL, Headers{}, nil)
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
func (h *NotificationHub) Registrations(ctx context.Context) (*Registrations, []byte, error) {
	rawResponse, err := h.exec(ctx, getMethod, h.generateAPIURL("registrations"), Headers{}, nil)
	if err != nil {
		return nil, rawResponse, err
	}
	res := &Registrations{}
	if err = xml.Unmarshal(rawResponse, &res); err != nil {
		return nil, rawResponse, err
	}
	res.normalize()
	return res, rawResponse, nil
}

// Register sends a device registration to the Azure hub
func (h *NotificationHub) Register(ctx context.Context, r Registration) (RegistrationResult, []byte, error) {
	var (
		regRes  = RegistrationResult{}
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
		return regRes, nil, errors.New("Notification format not implemented")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)

	if r.RegistrationID != "" {
		method = putMethod
		regURL.Path = path.Join(regURL.Path, r.RegistrationID)
	}

	res, err := h.exec(ctx, method, regURL, headers, bytes.NewBufferString(payload))

	if err == nil {
		if err = xml.Unmarshal(res, &regRes); err != nil {
			return regRes, res, err
		}
		regRes.RegistrationContent.normalize()
	}
	return regRes, res, err
}

// RegisterWithTemplate sends a device registration with template to the Azure hub
func (h *NotificationHub) RegisterWithTemplate(ctx context.Context, r Registration) (RegistrationResult, []byte, error) {
	var (
		regRes  = RegistrationResult{}
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
		return regRes, nil, errors.New("Notification format not implemented")
	}
	payload = strings.Replace(payload, "{{Tags}}", r.Tags, 1)

	if r.RegistrationID != "" {
		method = putMethod
		regURL.Path = path.Join(regURL.Path, r.RegistrationID)
	}

	res, err := h.exec(ctx, method, regURL, headers, bytes.NewBufferString(payload))

	if err == nil {
		if err = xml.Unmarshal(res, &regRes); err != nil {
			return regRes, res, err
		}
		regRes.RegistrationContent.normalize()
	}
	return regRes, res, err
}
