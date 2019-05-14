package notificationhubs_test

import (
	"reflect"
	"testing"

	. "github.com/daresaydigital/azure-notificationhubs-go"
)

func TestNewNotification(t *testing.T) {
	var (
		testPayload = []byte("test payload")
		errfmt      = "NewNotification test case %d error. Expected %s: %v, got: %v"

		testCases = []struct {
			format               NotificationFormat
			payload              []byte
			expectedNotification *Notification
			hasErr               bool
		}{
			{
				format:               Template,
				payload:              testPayload,
				expectedNotification: &Notification{Format: Template, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               GcmFormat,
				payload:              testPayload,
				expectedNotification: &Notification{Format: GcmFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               AppleFormat,
				payload:              testPayload,
				expectedNotification: &Notification{Format: AppleFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               BaiduFormat,
				payload:              testPayload,
				expectedNotification: &Notification{Format: BaiduFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               KindleFormat,
				payload:              testPayload,
				expectedNotification: &Notification{Format: KindleFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               WindowsFormat,
				payload:              testPayload,
				expectedNotification: &Notification{Format: WindowsFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               WindowsPhoneFormat,
				payload:              testPayload,
				expectedNotification: &Notification{Format: WindowsPhoneFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               NotificationFormat("unknown_format"),
				payload:              testPayload,
				expectedNotification: nil,
				hasErr:               true,
			},
		}
	)

	for i, testCase := range testCases {
		obtainedNotification, obtainedErr := NewNotification(testCase.format, testCase.payload)

		if !reflect.DeepEqual(obtainedNotification, testCase.expectedNotification) {
			t.Errorf(errfmt, i, "Notification", testCase.expectedNotification, obtainedNotification)
		}

		if (obtainedErr != nil) != testCase.hasErr {
			t.Errorf(errfmt, i, "hasError", testCase.hasErr, obtainedErr != nil)
		}
	}
}
