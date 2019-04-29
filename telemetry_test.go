package notificationhubs

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewNotificationTelemetryFromLocationURL(t *testing.T) {
	tests := []struct {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotificationTelemetryFromLocationURL(tt.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotificationTelemetryFromLocationURL() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNewNotificationTelemetryFromHTTPResponse(t *testing.T) {
	type args struct {
		response *http.Response
	}
	tests := []struct {
		name string
		args args
		want *NotificationTelemetry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotificationTelemetryFromHTTPResponse(tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotificationTelemetryFromHTTPResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
