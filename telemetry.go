package notificationhubs

import (
	"context"
	"encoding/xml"
	"path"
)

// NotificationDetails reads one specific registration
func (h *NotificationHub) NotificationDetails(ctx context.Context, deviceID string) (*RegistrationResult, []byte, error) {
	var (
		res    = &RegistrationResult{}
		regURL = h.generateAPIURL("registrations")
	)
	regURL.Path = path.Join(regURL.Path, deviceID)
	rawResponse, err := h.exec(ctx, getMethod, regURL, Headers{}, nil)
	if err != nil {
		return nil, rawResponse, err
	}
	if err = xml.Unmarshal(rawResponse, &res); err != nil {
		return nil, rawResponse, err
	}
	res.RegistrationContent.normalize()
	return res, rawResponse, nil
}
