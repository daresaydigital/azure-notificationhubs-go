package notificationhubs_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/daresaydigital/azure-notificationhubs-go/notificationhubs"
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

func Test_NotificationHubendFanout(t *testing.T) {
	nhub, notification, mockClient := initTestItems()

	mockClient.execFunc = func(obtainedReq *http.Request) ([]byte, error) {
		gotURL := obtainedReq.URL.String()
		if gotURL != messagesURL {
			t.Errorf(errfmt, "request URL", messagesURL, gotURL)
		}

		if obtainedReq.Method != "POST" {
			t.Errorf(errfmt, "request Method", "POST", obtainedReq.Method)
		}

		b, _ := ioutil.ReadAll(obtainedReq.Body)
		if string(b) != string(notification.Payload) {
			t.Errorf(errfmt, "request Body", string(notification.Payload), b)
		}

		if obtainedReq.Header.Get("Content-Type") != notification.Format.GetContentType() {
			t.Errorf(errfmt, "Content-Type header", notification.Format.GetContentType(), obtainedReq.Header.Get("Content-Type"))
		}

		if obtainedReq.Header.Get("ServiceBusNotification-Format") != string(notification.Format) {
			t.Errorf(errfmt, "ServiceBusNotification-Format header", notification.Format, obtainedReq.Header.Get("ServiceBusNotification-Format"))
		}

		if obtainedReq.Header.Get("ServiceBusNotification-Tags") != "" {
			t.Errorf(errfmt, "ServiceBusNotification-Tags", "", obtainedReq.Header.Get("ServiceBusNotification-Tags"))
		}

		obtainedAuthToken := obtainedReq.Header.Get("Authorization")
		expectedTokenPrefix := "SharedAccessSignature "

		var authTokenParams string
		if !strings.HasPrefix(obtainedAuthToken, expectedTokenPrefix) {
			t.Fatalf(errfmt, "auth token prefix", expectedTokenPrefix, strings.Split(obtainedAuthToken, " ")[0])
		} else {
			authTokenParams = strings.TrimPrefix(obtainedAuthToken, expectedTokenPrefix)
		}

		queryMap, _ := url.ParseQuery(authTokenParams)

		if len(queryMap["sr"]) == 0 || queryMap["sr"][0] != sasURIString {
			t.Errorf(errfmt, "token target uri", sasURIString, queryMap["sr"])
		}

		expectedSig := "cy3Y21BlsAw8slr5TnmSM3pilYBC8T7k3oNqOUKvE2g="
		if len(queryMap["sig"]) == 0 || queryMap["sig"][0] != expectedSig {
			t.Errorf(errfmt, "token signature", expectedSig, queryMap["sig"][0])
		}

		obtainedExpStr := queryMap["se"]
		if len(obtainedExpStr) == 0 {
			t.Errorf(errfmt, "token expiration", mockTimeGeneratorFunc.GenerateTimestamp(), obtainedExpStr)
		}

		obtainedExp, err := strconv.Atoi(obtainedExpStr[0])
		if err != nil || int64(obtainedExp) != mockTimeGeneratorFunc.GenerateTimestamp() {
			t.Errorf(errfmt, "token expiration", mockTimeGeneratorFunc.GenerateTimestamp(), obtainedExp)
		}

		if len(queryMap["skn"]) == 0 || queryMap["skn"][0] != nhub.SasKeyName {
			t.Errorf(errfmt, "token sas key name", nhub.SasKeyName, queryMap["skn"])
		}

		return nil, nil
	}

	b, err := nhub.Send(context.Background(), notification, nil)
	if b != nil {
		t.Errorf(errfmt, "byte", nil, b)
	}

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
}

func Test_NotificationHubendCategories(t *testing.T) {
	var (
		orTags                         = []string{"tag1", "tag2"}
		nhub, notification, mockClient = initTestItems()
	)

	mockClient.execFunc = func(obtainedReq *http.Request) ([]byte, error) {
		expectedTags := strings.Join(orTags, " || ")
		if obtainedReq.Header.Get("ServiceBusNotification-Tags") != expectedTags {
			t.Errorf(errfmt, "ServiceBusNotification-Tags", expectedTags, obtainedReq.Header.Get("ServiceBusNotification-Tags"))
		}

		gotURL := obtainedReq.URL.String()
		if gotURL != messagesURL {
			t.Errorf(errfmt, "URL", messagesURL, gotURL)
		}

		return nil, nil
	}

	b, err := nhub.Send(context.Background(), notification, orTags)
	if b != nil {
		t.Errorf(errfmt, "byte", nil, b)
	}

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
}

func Test_NotificationSendError(t *testing.T) {
	var (
		expectedError                  = errors.New("test error")
		nhub, notification, mockClient = initTestItems()
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, error) {
		if reqURL := req.URL.String(); reqURL != messagesURL {
			t.Errorf(errfmt, "URL", messagesURL, reqURL)
		}
		return nil, expectedError
	}

	b, obtainedErr := nhub.Send(context.Background(), notification, nil)
	if b != nil {
		t.Errorf(errfmt, "Send []byte", nil, b)
	}
	if !strings.Contains(obtainedErr.Error(), expectedError.Error()) {
		t.Errorf(errfmt, "Send error", expectedError, obtainedErr)
	}
}

func Test_NotificationScheduleSuccess(t *testing.T) {
	nhub, notification, mockClient := initTestItems()

	mockClient.execFunc = func(obtainedReq *http.Request) ([]byte, error) {
		gotURL := obtainedReq.URL.String()
		if gotURL != schedulesURL {
			t.Errorf(errfmt, "URL", schedulesURL, gotURL)
		}

		return nil, nil
	}

	b, err := nhub.Schedule(context.Background(), notification, nil, time.Now().Add(time.Minute))
	if b != nil {
		t.Errorf(errfmt, "byte", nil, b)
	}

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
}

func Test_NotificationScheduleOutdated(t *testing.T) {
	var (
		expectedError         = errors.New("You can not schedule a notification in the past")
		nhub, notification, _ = initTestItems()
	)
	b, err := nhub.Schedule(context.Background(), notification, nil, time.Now().Add(-time.Minute))
	if b != nil {
		t.Errorf(errfmt, "byte", nil, b)
	}

	if !strings.Contains(err.Error(), expectedError.Error()) {
		t.Errorf(errfmt, "Send error", expectedError, err)
	}
}

func Test_NotificationScheduleError(t *testing.T) {
	var (
		expectedError                  = errors.New("test schedule error")
		nhub, notification, mockClient = initTestItems()
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, error) {
		gotURL := req.URL.String()
		if gotURL != schedulesURL {
			t.Errorf(errfmt, "URL", schedulesURL, gotURL)
		}

		return nil, expectedError
	}

	b, obtainedErr := nhub.Schedule(context.Background(), notification, nil, time.Now().Add(time.Minute))
	if b != nil {
		t.Errorf(errfmt, "Send []byte", nil, b)
	}

	if !strings.Contains(obtainedErr.Error(), expectedError.Error()) {
		t.Errorf(errfmt, "Send error", expectedError, obtainedErr)
	}
}
