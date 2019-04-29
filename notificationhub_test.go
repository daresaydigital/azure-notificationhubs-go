package notificationhubs_test

import (
	"net/url"
	"testing"

	. "github.com/daresaydigital/azure-notificationhubs-go"
)

func Test_NewNotificationHub(t *testing.T) {
	var (
		errfmt      = "NewNotificationHub test case %d error. Expected %s: %v, got: %v"
		queryString = url.Values{apiVersionParam: {apiVersionValue}}.Encode()
		testCases   = []struct {
			connectionString string
			expectedHub      *mockNotificationHub
		}{
			{
				connectionString: connectionString,
				expectedHub: &mockNotificationHub{
					SasKeyValue: "testAccessKey",
					SasKeyName:  "testAccessKeyName",
					HubURL:      &url.URL{Host: "testhub-ns.servicebus.windows.net", Scheme: defaultScheme, Path: hubPath, RawQuery: queryString},
				},
			},
			{
				connectionString: "wrong_connection_string",
				expectedHub: &mockNotificationHub{
					SasKeyValue: "",
					SasKeyName:  "",
					HubURL:      &url.URL{Host: "", Scheme: defaultScheme, Path: hubPath, RawQuery: queryString},
				},
			},
		}
	)

	for i, testCase := range testCases {
		obtainedNotificationHub := NewNotificationHub(testCase.connectionString, hubPath)

		if obtainedNotificationHub.SasKeyValue != testCase.expectedHub.SasKeyValue {
			t.Errorf(errfmt, i, "NotificationHub.SasKeyValue", testCase.expectedHub.SasKeyValue, obtainedNotificationHub.SasKeyValue)
		}

		if obtainedNotificationHub.SasKeyName != testCase.expectedHub.SasKeyName {
			t.Errorf(errfmt, i, "NotificationHub.SasKeyName", testCase.expectedHub.SasKeyName, obtainedNotificationHub.SasKeyName)
		}

		wantURL := testCase.expectedHub.HubURL.String()
		gotURL := obtainedNotificationHub.HubURL.String()
		if gotURL != wantURL {
			t.Errorf(errfmt, i, "NotificationHub.hubURL", wantURL, gotURL)
		}
	}
}
