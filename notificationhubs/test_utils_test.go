package notificationhubs_test

import (
	"net/http"
	"net/url"
	"time"

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
	errfmt           = "Expected %s: \n%v\ngot:\n%v"
	postMethod       = "POST"
	putMethod        = "PUT"
	getMethod        = "GET"
)

var (
	endOfEpoch, _         = time.Parse("2006-01-02T15:04:05.000Z", "9999-12-31T23:59:59.999Z")
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

func initNotificationTestItems() (*NotificationHub, *Notification, *mockHubHTTPClient) {
	var (
		notification, _  = NewNotification(Template, []byte("test payload"))
		nhub, mockClient = initTestItems()
	)
	return nhub, notification, mockClient
}

func initTestItems() (*NotificationHub, *mockHubHTTPClient) {
	var (
		mockClient = &mockHubHTTPClient{}
		nhub       = NewNotificationHub(connectionString, hubPath)
	)
	nhub.SetHTTPClient(mockClient)
	nhub.SetExpirationTimeGenerator(mockTimeGeneratorFunc)
	return nhub, mockClient
}
