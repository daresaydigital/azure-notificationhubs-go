package notificationhubs_test

import (
	"testing"

	. "github.com/daresaydigital/azure-notificationhubs-go/notificationhubs"
)

func TestNotificationFormat_GetContentType(t *testing.T) {
	var (
		testCases = []struct {
			format   NotificationFormat
			expected string
		}{
			{
				format:   Template,
				expected: "application/json",
			},
			{
				format:   GcmFormat,
				expected: "application/json",
			},
			{
				format:   AppleFormat,
				expected: "application/json",
			},
			{
				format:   BaiduFormat,
				expected: "application/json",
			},
			{
				format:   KindleFormat,
				expected: "application/json",
			},
			{
				format:   WindowsFormat,
				expected: "application/xml",
			},
			{
				format:   WindowsPhoneFormat,
				expected: "application/xml",
			},
		}
	)

	for _, testCase := range testCases {
		obtained := testCase.format.GetContentType()
		if obtained != testCase.expected {
			t.Errorf("NotificationFormat '%s' GetContentType(). Expected '%s', got '%s'", testCase.format, testCase.expected, obtained)
		}
	}
}

func TestNotificationFormat_IsValid(t *testing.T) {
	var (
		testCases = []struct {
			format  NotificationFormat
			isValid bool
		}{
			{
				format:  Template,
				isValid: true,
			},
			{
				format:  GcmFormat,
				isValid: true,
			},
			{
				format:  AppleFormat,
				isValid: true,
			},
			{
				format:  BaiduFormat,
				isValid: true,
			},
			{
				format:  KindleFormat,
				isValid: true,
			},
			{
				format:  WindowsFormat,
				isValid: true,
			},
			{
				format:  WindowsPhoneFormat,
				isValid: true,
			},
			{
				format:  NotificationFormat("wrong_format"),
				isValid: false,
			},
		}
	)

	for _, testCase := range testCases {
		obtained := testCase.format.IsValid()
		if obtained != testCase.isValid {
			t.Errorf("NotificationFormat '%s' isValid(). Expected '%t', got '%t'", testCase.format, testCase.isValid, obtained)
		}
	}
}
