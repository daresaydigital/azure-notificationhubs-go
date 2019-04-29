package notificationhubs_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	. "github.com/daresaydigital/azure-notificationhubs-go"
)

func Test_RegisterApple(t *testing.T) {
	var (
		nhub, mockClient = initTestItems()
		registration     = Registration{
			Tags:               "tag1,tag2,tag3",
			DeviceID:           "ABCDEFG",
			NotificationFormat: AppleFormat,
		}
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, error) {
		gotMethod := req.Method
		if gotMethod != postMethod {
			t.Errorf(errfmt, "method", postMethod, gotMethod)
		}
		gotURL := req.URL.String()
		if gotURL != registrationsURL {
			t.Errorf(errfmt, "URL", registrationsURL, gotURL)
		}
		data, e := ioutil.ReadFile("./fixtures/appleRegistrationResult.xml")
		if e != nil {
			return nil, e
		}
		return data, nil
	}

	result, data, err := nhub.Register(context.Background(), registration)

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
	if data == nil {
		t.Errorf("Register response empty")
	} else {
		var (
			publishedTime, _ = time.Parse("2006-01-02T15:04:05Z", "2019-04-20T09:10:11Z")
			updatedTime, _   = time.Parse("2006-01-02T15:04:05Z", "2019-04-23T09:10:11Z")
		)
		expectedResult := RegistrationResult{
			ID:        "https://testhub-ns.servicebus.windows.net/testhub/registrations/8247220326459738692-7748251457295609952-3?api-version=2015-01",
			Title:     "8247220326459738692-7748251457295609952-3",
			Published: publishedTime,
			Updated:   updatedTime,
			RegistrationContent: &RegistrationContent{
				AppleRegistrationDescription: nil,
				GcmRegistrationDescription:   nil,
				RegistratedDevice: &RegistratedDevice{
					ETag:           "1",
					ExpirationTime: endOfEpoch,
					RegistrationID: "8247220326459738692-7748251457295609952-3",
					Tags:           []string{"tag1", "tag2", "tag3"},
					DeviceID:       "ABCDEFG",
				},
				Format: AppleFormat,
			},
		}
		if expectedResult.ID != result.ID {
			t.Errorf(errfmt, "", expectedResult.ID, result.ID)
		}
		if expectedResult.Title != result.Title {
			t.Errorf(errfmt, "", expectedResult.Title, result.Title)
		}
		if expectedResult.Published != result.Published {
			t.Errorf(errfmt, "", expectedResult.Published, result.Published)
		}
		if expectedResult.Updated != result.Updated {
			t.Errorf(errfmt, "", expectedResult.Updated, result.Updated)
		}
		if expectedResult.RegistrationContent.AppleRegistrationDescription != result.RegistrationContent.AppleRegistrationDescription {
			t.Errorf(errfmt, "", expectedResult.RegistrationContent.AppleRegistrationDescription, result.RegistrationContent.AppleRegistrationDescription)
		}
		if expectedResult.RegistrationContent.GcmRegistrationDescription != result.RegistrationContent.GcmRegistrationDescription {
			t.Errorf(errfmt, "", expectedResult.RegistrationContent.GcmRegistrationDescription, result.RegistrationContent.GcmRegistrationDescription)
		}
		if !reflect.DeepEqual(result.RegistrationContent.RegistratedDevice, expectedResult.RegistrationContent.RegistratedDevice) {
			t.Errorf(errfmt, "registration result", expectedResult.RegistrationContent.RegistratedDevice, result.RegistrationContent.RegistratedDevice)
		}
		if expectedResult.RegistrationContent.Format != result.RegistrationContent.Format {
			t.Errorf(errfmt, "", expectedResult.RegistrationContent.Format, result.RegistrationContent.Format)
		}
	}
}

func Test_RegisterGcm(t *testing.T) {
	var (
		nhub, mockClient = initTestItems()
		registration     = Registration{
			Tags:               "tag1,tag3",
			DeviceID:           "ANDROIDID",
			NotificationFormat: GcmFormat,
		}
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, error) {
		gotMethod := req.Method
		if gotMethod != postMethod {
			t.Errorf(errfmt, "method", postMethod, gotMethod)
		}
		gotURL := req.URL.String()
		if gotURL != registrationsURL {
			t.Errorf(errfmt, "URL", registrationsURL, gotURL)
		}
		data, e := ioutil.ReadFile("./fixtures/androidRegistrationResult.xml")
		if e != nil {
			return nil, e
		}
		return data, nil
	}

	result, data, err := nhub.Register(context.Background(), registration)

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
	if data == nil {
		t.Errorf("Register response empty")
	} else {
		var (
			publishedTime, _ = time.Parse("2006-01-02T15:04:05Z", "2019-04-20T09:19:06Z")
			updatedTime, _   = time.Parse("2006-01-02T15:04:05Z", "2019-04-23T09:19:06Z")
		)
		expectedResult := RegistrationResult{
			ID:        "https://testhub-ns.servicebus.windows.net/testhub/registrations/4603854756731398046-26535929789529194-1?api-version=2015-01",
			Title:     "4603854756731398046-26535929789529194-1",
			Published: publishedTime,
			Updated:   updatedTime,
			RegistrationContent: &RegistrationContent{
				AppleRegistrationDescription: nil,
				GcmRegistrationDescription:   nil,
				RegistratedDevice: &RegistratedDevice{
					ETag:              "1",
					ExpirationTime:    endOfEpoch,
					RegistrationID:    "4603854756731398046-26535929789529194-1",
					TagsString:        nil,
					Tags:              []string{"tag1", "tag3"},
					DeviceToken:       nil,
					GcmRegistrationID: nil,
					DeviceID:          "ANDROIDID",
				},
				Format: GcmFormat,
			},
		}
		if expectedResult.ID != result.ID {
			t.Errorf(errfmt, "", expectedResult.ID, result.ID)
		}
		if expectedResult.Title != result.Title {
			t.Errorf(errfmt, "", expectedResult.Title, result.Title)
		}
		if expectedResult.Published != result.Published {
			t.Errorf(errfmt, "", expectedResult.Published, result.Published)
		}
		if expectedResult.Updated != result.Updated {
			t.Errorf(errfmt, "", expectedResult.Updated, result.Updated)
		}
		if expectedResult.RegistrationContent.AppleRegistrationDescription != result.RegistrationContent.AppleRegistrationDescription {
			t.Errorf(errfmt, "", expectedResult.RegistrationContent.AppleRegistrationDescription, result.RegistrationContent.AppleRegistrationDescription)
		}
		if expectedResult.RegistrationContent.GcmRegistrationDescription != result.RegistrationContent.GcmRegistrationDescription {
			t.Errorf(errfmt, "", expectedResult.RegistrationContent.GcmRegistrationDescription, result.RegistrationContent.GcmRegistrationDescription)
		}
		if !reflect.DeepEqual(result.RegistrationContent.RegistratedDevice, expectedResult.RegistrationContent.RegistratedDevice) {
			t.Errorf(errfmt, "registration result", expectedResult.RegistrationContent.RegistratedDevice, result.RegistrationContent.RegistratedDevice)
		}
		if expectedResult.RegistrationContent.Format != result.RegistrationContent.Format {
			t.Errorf(errfmt, "", expectedResult.RegistrationContent.Format, result.RegistrationContent.Format)
		}
	}
}

func Test_Registrations(t *testing.T) {
	var (
		nhub, mockClient = initTestItems()
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, error) {
		gotMethod := req.Method
		if gotMethod != getMethod {
			t.Errorf(errfmt, "method", getMethod, gotMethod)
		}
		gotURL := req.URL.String()
		if gotURL != registrationsURL {
			t.Errorf(errfmt, "URL", registrationsURL, gotURL)
		}
		data, e := ioutil.ReadFile("./fixtures/registrationsResult.xml")
		if e != nil {
			return nil, e
		}
		return data, nil
	}

	result, data, err := nhub.Registrations(context.Background())

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
	if data == nil {
		t.Errorf("Registrations response empty")
	} else {
		if result.ID != "https://testhub-ns.servicebus.windows.net/testhub/registrations?api-version=2015-01" {
			t.Errorf(errfmt, "id", "https://testhub-ns.servicebus.windows.net/testhub/registrations?api-version=2015-01", result.ID)
		}
		if len(result.Entries) != 4 {
			t.Errorf(errfmt, "entries", 4, len(result.Entries))
		}
		if result.Entries[0].RegistrationContent.Format != AppleFormat {
			t.Errorf(errfmt, "device format", AppleFormat, result.Entries[0].RegistrationContent.Format)
		}
		if result.Entries[0].RegistrationContent.RegistratedDevice.DeviceID != "ABCDEF" {
			t.Errorf(errfmt, "device format", "ABCDEF", result.Entries[0].RegistrationContent.RegistratedDevice.DeviceID)
		}
		if result.Entries[1].RegistrationContent.Format != AppleFormat {
			t.Errorf(errfmt, "device format", AppleFormat, result.Entries[1].RegistrationContent.Format)
		}
		if result.Entries[1].RegistrationContent.RegistratedDevice.DeviceID != "QWERTY" {
			t.Errorf(errfmt, "device format", "QWERTY", result.Entries[1].RegistrationContent.RegistratedDevice.DeviceID)
		}
		if result.Entries[2].RegistrationContent.Format != AppleFormat {
			t.Errorf(errfmt, "device format", AppleFormat, result.Entries[2].RegistrationContent.Format)
		}
		if result.Entries[2].RegistrationContent.RegistratedDevice.DeviceID != "ZXCVBN" {
			t.Errorf(errfmt, "device format", "ZXCVBN", result.Entries[2].RegistrationContent.RegistratedDevice.DeviceID)
		}
		if result.Entries[3].RegistrationContent.Format != GcmFormat {
			t.Errorf(errfmt, "device format", GcmFormat, result.Entries[3].RegistrationContent.Format)
		}
		if result.Entries[3].RegistrationContent.RegistratedDevice.DeviceID != "ANDROIDID" {
			t.Errorf(errfmt, "device format", "ANDROIDID", result.Entries[3].RegistrationContent.RegistratedDevice.DeviceID)
		}
	}
}
