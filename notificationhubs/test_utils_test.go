package notificationhubs_test

import (
	"net/http"
	"net/url"

	. "github.com/daresaydigital/azure-notificationhubs-go/notificationhubs"
	"github.com/daresaydigital/azure-notificationhubs-go/notificationhubs/utils"
	"gopkg.in/xmlpath.v2"
)

// Internal constants for testing
const (
	connectionString = "Endpoint=sb://testhub-ns.servicebus.windows.net/;SharedAccessKeyName=testAccessKeyName;SharedAccessKey=testAccessKey"
	messagesURL      = "https://testhub-ns.servicebus.windows.net/testhub/messages?api-version=2015-01"
	schedulesURL     = "https://testhub-ns.servicebus.windows.net/testhub/schedulednotifications?api-version=2015-01"
	registrationsURL = "https://testhub-ns.servicebus.windows.net/testhub/registrations?api-version=2015-01"
	hubPath          = "testhub"
	apiVersionParam  = "api-version"
	apiVersionValue  = "2015-01" // Looks old but the API is the same
	directParam      = "direct"
	defaultScheme    = "https"
	errfmt           = "Expected %s: %v, got: %v"
)

var (
	mockTimeGeneratorFunc = utils.ExpirationTimeGeneratorFunc(func() int64 { return 123 })
	realTimeGeneratorFunc = utils.NewExpirationTimeGenerator()
	sasURIString          = (&url.URL{Host: "testhub-ns.servicebus.windows.net", Scheme: defaultScheme}).String()
)

type mockNotificationHub struct {
	SasKeyValue string
	SasKeyName  string
	HubURL      *url.URL

	client                  utils.HTTPClient
	expirationTimeGenerator utils.ExpirationTimeGenerator
	regIDPath               *xmlpath.Path
	eTagPath                *xmlpath.Path
	expTmPath               *xmlpath.Path
}

type mockHubHTTPClient struct {
	execFunc func(*http.Request) ([]byte, error)
}

func (mc *mockHubHTTPClient) Exec(req *http.Request) ([]byte, error) {
	return mc.execFunc(req)
}

func initTestItems() (*NotificationHub, *Notification, *mockHubHTTPClient) {
	var (
		notification, _ = NewNotification(Template, []byte("test payload"))
		mockClient      = &mockHubHTTPClient{}
		nhub            = NewNotificationHub(connectionString, hubPath)
	)
	nhub.SetHTTPClient(mockClient)
	nhub.SetExpirationTimeGenerator(mockTimeGeneratorFunc)
	return nhub, notification, mockClient
}
