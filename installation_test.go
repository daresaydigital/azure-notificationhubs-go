package notificationhubs_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	. "github.com/daresaydigital/azure-notificationhubs-go"
)

func Test_Install(t *testing.T) {
	var (
		nhub, mockClient = initTestItems()
		installation     = Installation{
			InstallationID: "0a92196c-20c3-4308-8046-c384c902d0ff",
			Tags:           []string{"tag1", "tag3"},
			PushChannel:    "ANDROIDID",
			Platform:       GCMPlatform,
		}
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, *http.Response, error) {
		gotMethod := req.Method
		if gotMethod != putMethod {
			t.Errorf(errfmt, "method", putMethod, gotMethod)
		}
		u, _ := url.Parse(installationsURL)
		u.Path += "/" + installation.InstallationID
		wantURL := u.String()
		gotURL := req.URL.String()
		if gotURL != wantURL {
			t.Errorf(errfmt, "URL", wantURL, gotURL)
		}
		return nil, nil, nil
	}

	err := nhub.Install(context.Background(), installation)

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
}

func Test_Installation(t *testing.T) {
	var (
		nhub, mockClient = initTestItems()
		installationID   = "0a92196c-20c3-4308-8046-c384c902d0ff"
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, *http.Response, error) {
		gotMethod := req.Method
		if gotMethod != getMethod {
			t.Errorf(errfmt, "method", putMethod, gotMethod)
		}
		u, _ := url.Parse(installationsURL)
		u.Path += "/" + installationID
		wantURL := u.String()
		gotURL := req.URL.String()
		if gotURL != wantURL {
			t.Errorf(errfmt, "URL", wantURL, gotURL)
		}
		data, e := ioutil.ReadFile("./fixtures/gcmInstallationResult.json")
		if e != nil {
			return nil, nil, e
		}
		return data, nil, nil
	}

	data, result, err := nhub.Installation(context.Background(), installationID)

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
	if data == nil {
		t.Errorf("Install response empty")
	} else {
		expectedResult := &Installation{
			ExpirationTime: &endOfEpoch,
			InstallationID: "0a92196c-20c3-4308-8046-c384c902d0ff",
			Tags:           []string{"tag1", "tag3"},
			PushChannel:    "ANDROIDID",
			Platform:       GCMPlatform,
		}
		if !reflect.DeepEqual(result, expectedResult) {
			t.Errorf(errfmt, "installation result", result, expectedResult)
		}
	}
}

func Test_Update(t *testing.T) {
	var (
		nhub, mockClient = initTestItems()
		installationID   = "0a92196c-20c3-4308-8046-c384c902d0ff"
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, *http.Response, error) {
		gotMethod := req.Method
		if gotMethod != patchMethod {
			t.Errorf(errfmt, "method", putMethod, gotMethod)
		}
		u, _ := url.Parse(installationsURL)
		u.Path += "/" + installationID
		wantURL := u.String()
		gotURL := req.URL.String()
		if gotURL != wantURL {
			t.Errorf(errfmt, "URL", wantURL, gotURL)
		}
		return nil, nil, nil
	}

	err := nhub.Update(context.Background(), installationID,
		SetPushChannel("pushChannel"),
		SetTags("tag1", "tag2"),
		AddTag("tag"),
		RemoveTag("tag"),
		AddTemplate("name", InstallationTemplate{}),
		SetTemplateBody("name", "body"),
		SetTemplateHeaders("name", map[string]string{}),
		SetTemplateTags("name", "tag1", "tag2"),
		AddTemplateTag("name", "tag"),
		RemoveTemplateTag("name", "tag"),
		RemoveTemplate("name"),
		SetSecondaryTilePushChannel("name", "pushChannel"),
		SetSecondaryTileTags("name", "tag1", "tag2"),
		AddSecondaryTileTag("name", "tag"),
		RemoveSecondaryTileTag("name", "tag"),
		SetSecondaryTileTemplates("name", map[string]InstallationTemplate{}),
		AddSecondaryTileTemplate("name", "templateName", InstallationTemplate{}),
		SetSecondaryTileTemplateBody("name", "template", "body"),
		SetSecondaryTileTemplateHeaders("name", "template", map[string]string{}),
		SetSecondaryTileTemplateTags("name", "template", "tag1", "tag2"),
		RemoveSecondaryTileTemplate("name", "template"),
	)

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
}

func Test_Uninstall(t *testing.T) {
	var (
		nhub, mockClient = initTestItems()
		installationID   = "0a92196c-20c3-4308-8046-c384c902d0ff"
	)

	mockClient.execFunc = func(req *http.Request) ([]byte, *http.Response, error) {
		gotMethod := req.Method
		if gotMethod != deleteMethod {
			t.Errorf(errfmt, "method", deleteMethod, gotMethod)
		}
		u, _ := url.Parse(installationsURL)
		u.Path += "/" + installationID
		wantURL := u.String()
		gotURL := req.URL.String()
		if gotURL != wantURL {
			t.Errorf(errfmt, "URL", wantURL, gotURL)
		}
		return nil, nil, nil
	}

	err := nhub.Uninstall(context.Background(), installationID)

	if err != nil {
		t.Errorf(errfmt, "error", nil, err)
	}
}
