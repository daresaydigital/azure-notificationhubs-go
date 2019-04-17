package notificationhubs_test

import (
	"reflect"
	"testing"

	nh "github.com/daresaydigital/azure-notificationhubs-go/notificationhubs"
)

func TestNewNotification(t *testing.T) {
	var (
		testPayload = []byte("test payload")
		errfmt      = "NewNotification test case %d error. Expected %s: %v, got: %v"

		testCases = []struct {
			format               nh.NotificationFormat
			payload              []byte
			expectedNotification *nh.Notification
			hasErr               bool
		}{
			{
				format:               nh.Template,
				payload:              testPayload,
				expectedNotification: &nh.Notification{Format: nh.Template, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               nh.AndroidFormat,
				payload:              testPayload,
				expectedNotification: &nh.Notification{Format: nh.AndroidFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               nh.AppleFormat,
				payload:              testPayload,
				expectedNotification: &nh.Notification{Format: nh.AppleFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               nh.BaiduFormat,
				payload:              testPayload,
				expectedNotification: &nh.Notification{Format: nh.BaiduFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               nh.KindleFormat,
				payload:              testPayload,
				expectedNotification: &nh.Notification{Format: nh.KindleFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               nh.WindowsFormat,
				payload:              testPayload,
				expectedNotification: &nh.Notification{Format: nh.WindowsFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               nh.WindowsPhoneFormat,
				payload:              testPayload,
				expectedNotification: &nh.Notification{Format: nh.WindowsPhoneFormat, Payload: testPayload},
				hasErr:               false,
			},
			{
				format:               nh.NotificationFormat("unknown_format"),
				payload:              testPayload,
				expectedNotification: nil,
				hasErr:               true,
			},
		}
	)

	for i, testCase := range testCases {
		obtainedNotification, obtainedErr := nh.NewNotification(testCase.format, testCase.payload)

		if !reflect.DeepEqual(obtainedNotification, testCase.expectedNotification) {
			t.Errorf(errfmt, i, "Notification", testCase.expectedNotification, obtainedNotification)
		}

		if (obtainedErr != nil) != testCase.hasErr {
			t.Errorf(errfmt, i, "hasError", testCase.hasErr, obtainedErr != nil)
		}
	}
}
