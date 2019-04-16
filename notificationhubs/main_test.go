package notificationhubs

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_NotificationFormatIsValid(t *testing.T) {
	testCases := []struct {
		format  NotificationFormat
		isValid bool
	}{
		{
			format:  Template,
			isValid: true,
		},
		{
			format:  AndroidFormat,
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

	for _, testCase := range testCases {
		obtained := testCase.format.IsValid()
		if obtained != testCase.isValid {
			t.Errorf("NotificationFormat '%s' isValid(). Expected '%t', got '%t'", testCase.format, testCase.isValid, obtained)
		}
	}
}

func Test_NotificationFormatGetContentType(t *testing.T) {
	testCases := []struct {
		format   NotificationFormat
		expected string
	}{
		{
			format:   Template,
			expected: "application/json",
		},
		{
			format:   AndroidFormat,
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

	for _, testCase := range testCases {
		obtained := testCase.format.GetContentType()
		if obtained != testCase.expected {
			t.Errorf("NotificationFormat '%s' GetContentType(). Expected '%s', got '%s'", testCase.format, testCase.expected, obtained)
		}
	}
}

func Test_NewNotication(t *testing.T) {
	testPayload := []byte("test payload")
	errfmt := "NewNotification test case %d error. Expected %s: %v, got: %v"

	testCases := []struct {
		format               NotificationFormat
		payload              []byte
		expectedNotification *Notification
		hasErr               bool
	}{
		{
			format:               Template,
			payload:              testPayload,
			expectedNotification: &Notification{Template, testPayload},
			hasErr:               false,
		},
		{
			format:               AndroidFormat,
			payload:              testPayload,
			expectedNotification: &Notification{AndroidFormat, testPayload},
			hasErr:               false,
		},
		{
			format:               AppleFormat,
			payload:              testPayload,
			expectedNotification: &Notification{AppleFormat, testPayload},
			hasErr:               false,
		},
		{
			format:               BaiduFormat,
			payload:              testPayload,
			expectedNotification: &Notification{BaiduFormat, testPayload},
			hasErr:               false,
		},
		{
			format:               KindleFormat,
			payload:              testPayload,
			expectedNotification: &Notification{KindleFormat, testPayload},
			hasErr:               false,
		},
		{
			format:               WindowsFormat,
			payload:              testPayload,
			expectedNotification: &Notification{WindowsFormat, testPayload},
			hasErr:               false,
		},
		{
			format:               WindowsPhoneFormat,
			payload:              testPayload,
			expectedNotification: &Notification{WindowsPhoneFormat, testPayload},
			hasErr:               false,
		},
		{
			format:               NotificationFormat("unknown_format"),
			payload:              testPayload,
			expectedNotification: nil,
			hasErr:               true,
		},
	}

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

func Test_NotificationString(t *testing.T) {
	n := &Notification{Template, []byte("test_payload")}

	expectedString := fmt.Sprintf("&{%s %s}", n.Format, n.Payload)
	obtainedString := n.String()
	if obtainedString != expectedString {
		t.Errorf("Expected: %s, got: %s", expectedString, obtainedString)
	}
}

func Test_NewNotificationHub(t *testing.T) {
	errfmt := "NewNotificationHub test case %d error. Expected %s: %v, got: %v"

	queryString := url.Values{apiVersionParam: {apiVersionValue}}.Encode()
	hubPath := "testhub"
	testCases := []struct {
		connectionString string
		expectedHub      *NotificationHub
	}{
		{
			connectionString: "Endpoint=sb://testhub-ns.servicebus.windows.net/;SharedAccessKeyName=testAccessKeyName;SharedAccessKey=testAccessKey",
			expectedHub: &NotificationHub{
				sasKeyValue:             "testAccessKey",
				sasKeyName:              "testAccessKeyName",
				hubURL:                  &url.URL{Host: "testhub-ns.servicebus.windows.net", Scheme: schemeDefault, Path: hubPath, RawQuery: queryString},
				client:                  &hubHTTPClient{&http.Client{}},
				expirationTimeGenerator: expirationTimeGeneratorFunc(generateExpirationTimestamp),
			},
		},
		{
			connectionString: "wrong_connection_string",
			expectedHub: &NotificationHub{
				sasKeyValue:             "",
				sasKeyName:              "",
				hubURL:                  &url.URL{Host: "", Scheme: schemeDefault, Path: hubPath, RawQuery: queryString},
				client:                  &hubHTTPClient{&http.Client{}},
				expirationTimeGenerator: expirationTimeGeneratorFunc(generateExpirationTimestamp),
			},
		},
	}

	for i, testCase := range testCases {
		obtainedNotificationHub := NewNotificationHub(testCase.connectionString, hubPath)

		if testCase.expectedHub.sasKeyValue != testCase.expectedHub.sasKeyValue {
			t.Errorf(errfmt, i, "NotificationHub.sasKeyValue", testCase.expectedHub.sasKeyValue, obtainedNotificationHub.sasKeyValue)
		}

		if obtainedNotificationHub.sasKeyName != testCase.expectedHub.sasKeyName {
			t.Errorf(errfmt, i, "NotificationHub.sasKeyName", testCase.expectedHub.sasKeyName, obtainedNotificationHub.sasKeyName)
		}

		wantURL := testCase.expectedHub.hubURL.String()
		gotURL := obtainedNotificationHub.hubURL.String()
		if gotURL != wantURL {
			t.Errorf(errfmt, i, "NotificationHub.hubURL", wantURL, gotURL)
		}

		expectedGeneratorType := reflect.ValueOf(testCase.expectedHub.expirationTimeGenerator).Type()
		obtainedGeneratorType := reflect.ValueOf(obtainedNotificationHub.expirationTimeGenerator).Type()
		if !obtainedGeneratorType.AssignableTo(expectedGeneratorType) {
			t.Errorf(errfmt, i, "NotificationHub.expirationTimeGenerator", expectedGeneratorType, obtainedGeneratorType)
		}
	}
}

type mockHubHTTPClient struct {
	execFunc func(*http.Request) ([]byte, error)
}

func (mc *mockHubHTTPClient) Exec(req *http.Request) ([]byte, error) {
	return mc.execFunc(req)
}

func Test_NotificationHubendFanout(t *testing.T) {
	var (
		errfmt       = "Expected %s: %v, got: %v"
		notification = &Notification{Template, []byte("test payload")}

		baseURL = &url.URL{
			Host:     "testHost",
			Scheme:   schemeDefault,
			Path:     "testPath",
			RawQuery: url.Values{"queryParam": {"queryValue"}}.Encode(),
		}
		sasURI = (&url.URL{Host: baseURL.Host, Scheme: baseURL.Scheme}).String()
	)

	mockClient := &mockHubHTTPClient{}

	nhub := &NotificationHub{
		sasKeyValue:             "testKeyValue",
		sasKeyName:              "testKeyName",
		hubURL:                  baseURL,
		client:                  mockClient,
		expirationTimeGenerator: expirationTimeGeneratorFunc(func() int64 { return 123 }),
	}

	msgURL := "https://testHost/testPath/messages?queryParam=queryValue"

	mockClient.execFunc = func(obtainedReq *http.Request) ([]byte, error) {
		gotURL := obtainedReq.URL.String()
		if gotURL != msgURL {
			t.Errorf(errfmt, "request URL", msgURL, gotURL)
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

		expectedURI := strings.ToLower(sasURI)
		if len(queryMap["sr"]) == 0 || queryMap["sr"][0] != expectedURI {
			t.Errorf(errfmt, "token target uri", expectedURI, queryMap["sr"])
		}

		expectedSig := "gbQ5tD5dkCLLu6FavSBKu4b7EAPeFqF7XEoDOada6ww="
		if len(queryMap["sig"]) == 0 || queryMap["sig"][0] != expectedSig {
			t.Errorf(errfmt, "token signature", expectedSig, queryMap["sig"][0])
		}

		obtainedExpStr := queryMap["se"]
		if len(obtainedExpStr) == 0 {
			t.Errorf(errfmt, "token expiration", nhub.expirationTimeGenerator.GenerateTimestamp(), obtainedExpStr)
		}

		obtainedExp, err := strconv.Atoi(obtainedExpStr[0])
		if err != nil || int64(obtainedExp) != nhub.expirationTimeGenerator.GenerateTimestamp() {
			t.Errorf(errfmt, "token expiration", nhub.expirationTimeGenerator.GenerateTimestamp(), obtainedExp)
		}

		if len(queryMap["skn"]) == 0 || queryMap["skn"][0] != nhub.sasKeyName {
			t.Errorf(errfmt, "token sas key name", nhub.sasKeyName, queryMap["skn"])
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
		errfmt = "Expected %s: %v, got: %v"

		orTags       = []string{"tag1", "tag2"}
		notification = &Notification{Template, []byte("test_payload")}

		baseURL = &url.URL{
			Host:     "testHost",
			Scheme:   schemeDefault,
			Path:     "testPath",
			RawQuery: url.Values{"queryParam": {"queryValue"}}.Encode(),
		}
	)

	mockClient := &mockHubHTTPClient{}

	nhub := &NotificationHub{
		sasKeyName:              "testKeyName",
		sasKeyValue:             "testKeyValue",
		hubURL:                  baseURL,
		client:                  mockClient,
		expirationTimeGenerator: expirationTimeGeneratorFunc(func() int64 { return 123 }),
	}

	msgURL := "https://testHost/testPath/messages?queryParam=queryValue"

	mockClient.execFunc = func(obtainedReq *http.Request) ([]byte, error) {
		expectedTags := strings.Join(orTags, " || ")
		if obtainedReq.Header.Get("ServiceBusNotification-Tags") != expectedTags {
			t.Errorf(errfmt, "ServiceBusNotification-Tags", expectedTags, obtainedReq.Header.Get("ServiceBusNotification-Tags"))
		}

		gotURL := obtainedReq.URL.String()
		if gotURL != msgURL {
			t.Errorf(errfmt, "URL", msgURL, gotURL)
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
		errfmt        = "Expected %s: %v, got: %v"
		expectedError = errors.New("test error")

		baseURL = &url.URL{
			Host:     "testHost",
			Scheme:   schemeDefault,
			Path:     "testPath",
			RawQuery: url.Values{"queryParam": {"queryValue"}}.Encode(),
		}
	)

	msgURL := "https://testHost/testPath/messages?queryParam=queryValue"

	mockClient := &mockHubHTTPClient{}
	mockClient.execFunc = func(req *http.Request) ([]byte, error) {
		if reqURL := req.URL.String(); reqURL != msgURL {
			t.Errorf(errfmt, "URL", msgURL, reqURL)
		}

		return nil, expectedError
	}

	nhub := &NotificationHub{
		sasKeyValue:             "testKeyValue",
		sasKeyName:              "testKeyName",
		hubURL:                  baseURL,
		client:                  mockClient,
		expirationTimeGenerator: expirationTimeGeneratorFunc(func() int64 { return 123 }),
	}

	b, obtainedErr := nhub.Send(context.Background(), &Notification{AndroidFormat, []byte("test payload")}, nil)
	if b != nil {
		t.Errorf(errfmt, "Send []byte", nil, b)
	}

	if !strings.Contains(obtainedErr.Error(), expectedError.Error()) {
		t.Errorf(errfmt, "Send error", expectedError, obtainedErr)
	}
}

func Test_NotificationScheduleSuccess(t *testing.T) {
	var (
		errfmt       = "Expected %s: %v, got: %v"
		notification = &Notification{Template, []byte("test_payload")}
		baseURL      = &url.URL{
			Host:     "testHost",
			Scheme:   schemeDefault,
			Path:     "testPath",
			RawQuery: url.Values{"queryParam": {"queryValue"}}.Encode(),
		}
	)

	mockClient := &mockHubHTTPClient{}

	nhub := &NotificationHub{
		sasKeyValue:             "testKeyValue",
		sasKeyName:              "testKeyName",
		hubURL:                  baseURL,
		client:                  mockClient,
		expirationTimeGenerator: expirationTimeGeneratorFunc(func() int64 { return 123 }),
	}

	schURL := "https://testHost/testPath/schedulednotifications?queryParam=queryValue"

	mockClient.execFunc = func(obtainedReq *http.Request) ([]byte, error) {
		gotURL := obtainedReq.URL.String()
		if gotURL != schURL {
			t.Errorf(errfmt, "URL", schURL, gotURL)
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
		errfmt       = "Expected %s: %v, got: %v"
		notification = &Notification{Template, []byte("test_payload")}

		baseURL = &url.URL{
			Host:     "testHost",
			Scheme:   schemeDefault,
			Path:     "testPath",
			RawQuery: url.Values{"queryParam": {"queryValue"}}.Encode(),
		}
	)

	mockClient := &mockHubHTTPClient{}

	nhub := &NotificationHub{
		sasKeyValue:             "testKeyValue",
		sasKeyName:              "testKeyName",
		hubURL:                  baseURL,
		client:                  mockClient,
		expirationTimeGenerator: expirationTimeGeneratorFunc(func() int64 { return 123 }),
	}

	schURL := "https://testHost/testPath/messages?queryParam=queryValue"
	mockClient.execFunc = func(obtainedReq *http.Request) ([]byte, error) {
		gotURL := obtainedReq.URL.String()
		if gotURL != schURL {
			t.Errorf(errfmt, "URL", schURL, gotURL)
		}

		return nil, nil
	}

	b, err := nhub.Schedule(context.Background(), notification, nil, time.Now().Add(-time.Minute))
	if b != nil {
		t.Errorf(errfmt, "byte", nil, b)
	}

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
}

func Test_NotificationScheduleError(t *testing.T) {
	var (
		errfmt        = "Expected %s: %v, got: %v"
		expectedError = errors.New("test schedule error")

		baseURL = &url.URL{
			Host:     "testHost",
			Scheme:   schemeDefault,
			Path:     "testPath",
			RawQuery: url.Values{"queryParam": {"queryValue"}}.Encode(),
		}
	)

	schURL := "https://testHost/testPath/schedulednotifications?queryParam=queryValue"

	mockClient := &mockHubHTTPClient{}
	mockClient.execFunc = func(req *http.Request) ([]byte, error) {
		gotURL := req.URL.String()
		if gotURL != schURL {
			t.Errorf(errfmt, "URL", schURL, gotURL)
		}

		return nil, expectedError
	}

	nhub := &NotificationHub{
		sasKeyValue:             "testKeyValue",
		sasKeyName:              "testKeyName",
		hubURL:                  baseURL,
		client:                  mockClient,
		expirationTimeGenerator: expirationTimeGeneratorFunc(func() int64 { return 123 }),
	}

	b, obtainedErr := nhub.Schedule(context.Background(), &Notification{AndroidFormat, []byte("test payload")}, nil, time.Now().Add(time.Minute))
	if b != nil {
		t.Errorf(errfmt, "Send []byte", nil, b)
	}

	if !strings.Contains(obtainedErr.Error(), expectedError.Error()) {
		t.Errorf(errfmt, "Send error", expectedError, obtainedErr)
	}
}
