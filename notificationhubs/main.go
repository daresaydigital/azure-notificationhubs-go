// Package notificationhubs represents an http client
// for Microsoft Azure Notification Hub
// Originally a fork of Gozure https://github.com/Onefootball/gozure
package notificationhubs

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/xmlpath.v2"
)

// NewNotificationHub initializes and returns NotificationHub pointer
func NewNotificationHub(connectionString, hubPath string) *NotificationHub {
	connData := strings.Split(connectionString, ";")

	hub := &NotificationHub{
		hubURL: &url.URL{},
	}

	for _, connItem := range connData {
		if strings.HasPrefix(connItem, paramEndpoint) {
			hubURL, err := url.Parse(connItem[len(paramEndpoint):])
			if err == nil {
				hub.hubURL = hubURL
			}
			continue
		}

		if strings.HasPrefix(connItem, paramSaasKeyName) {
			hub.sasKeyName = connItem[len(paramSaasKeyName):]
			continue
		}

		if strings.HasPrefix(connItem, paramSaasKeyValue) {
			hub.sasKeyValue = connItem[len(paramSaasKeyValue):]
			continue
		}
	}

	if hub.hubURL.Scheme == schemeServiceBus || hub.hubURL.Scheme == "" {
		hub.hubURL.Scheme = schemeDefault
	}

	hub.hubURL.Path = hubPath
	hub.hubURL.RawQuery = url.Values{apiVersionParam: {apiVersionValue}}.Encode()

	hub.client = &hubHTTPClient{&http.Client{}}
	hub.expirationTimeGenerator = expirationTimeGeneratorFunc(generateExpirationTimestamp)

	hub.regIDPath = xmlpath.MustCompile("/entry/content/*/RegistrationId")
	hub.eTagPath = xmlpath.MustCompile("/entry/content/*/ETag")
	hub.expTmPath = xmlpath.MustCompile("/entry/content/*/ExpirationTime")

	return hub
}

// NewNotification initalizes and returns Notification pointer
func NewNotification(format NotificationFormat, payload []byte) (*Notification, error) {
	if !format.IsValid() {
		return nil, fmt.Errorf("unknown format '%s'", format)
	}

	return &Notification{format, payload}, nil
}
