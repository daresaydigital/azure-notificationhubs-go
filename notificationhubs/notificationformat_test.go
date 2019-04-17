package notificationhubs_test

import (
	"testing"

	nh "github.com/daresaydigital/azure-notificationhubs-go/notificationhubs"
)

func TestNotificationFormat_GetContentType(t *testing.T) {
	var (
		testCases = []struct {
			format   nh.NotificationFormat
			expected string
		}{
			{
				format:   nh.Template,
				expected: "application/json",
			},
			{
				format:   nh.AndroidFormat,
				expected: "application/json",
			},
			{
				format:   nh.AppleFormat,
				expected: "application/json",
			},
			{
				format:   nh.BaiduFormat,
				expected: "application/json",
			},
			{
				format:   nh.KindleFormat,
				expected: "application/json",
			},
			{
				format:   nh.WindowsFormat,
				expected: "application/xml",
			},
			{
				format:   nh.WindowsPhoneFormat,
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
			format  nh.NotificationFormat
			isValid bool
		}{
			{
				format:  nh.Template,
				isValid: true,
			},
			{
				format:  nh.AndroidFormat,
				isValid: true,
			},
			{
				format:  nh.AppleFormat,
				isValid: true,
			},
			{
				format:  nh.BaiduFormat,
				isValid: true,
			},
			{
				format:  nh.KindleFormat,
				isValid: true,
			},
			{
				format:  nh.WindowsFormat,
				isValid: true,
			},
			{
				format:  nh.WindowsPhoneFormat,
				isValid: true,
			},
			{
				format:  nh.NotificationFormat("wrong_format"),
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
