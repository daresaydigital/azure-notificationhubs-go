package notificationhubs

import (
	"net/http"
	"reflect"
	"testing"
)

var tests = []struct {
	name string
	url  string
	want *NotificationTelemetry
}{
	{
		name: "Simple test",
		url:  "https://test-ns.servicebus.windows.net/testhub/messages/ABCDEFGH?api-version=2015-04",
		want: &NotificationTelemetry{
			NotificationMessageID: "ABCDEFGH",
		},
	},
	{
		name: "More advanced test",
		url:  "https://messages.servicebus.windows.net/messagebus/messages/3288835312934927344-986564390439048203-1?api-version=2016-10",
		want: &NotificationTelemetry{
			NotificationMessageID: "3288835312934927344-986564390439048203-1",
		},
	},
}

func TestNewNotificationTelemetryFromLocationURL(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotificationTelemetryFromLocationURL(tt.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotificationTelemetryFromLocationURL() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNewNotificationTelemetryFromHTTPResponse(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			header.Add("Location", tt.url)
			response := &http.Response{
				Header: header,
			}
			if got, _ := NewNotificationTelemetryFromHTTPResponse(response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotificationTelemetryFromHTTPResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
