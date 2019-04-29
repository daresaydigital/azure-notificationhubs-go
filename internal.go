package notificationhubs

// Internal constants
const (
	apiVersionParam          = "api-version"
	apiVersionValue          = "2015-01"
	telemetryAPIVersionValue = "2016-07"
	directParam              = "direct"

	// for connection string parsing
	schemeServiceBus  = "sb"
	schemeDefault     = "https"
	paramEndpoint     = "Endpoint="
	paramSaasKeyName  = "SharedAccessKeyName="
	paramSaasKeyValue = "SharedAccessKey="

	// Http methods
	deleteMethod = "DELETE"
	getMethod    = "GET"
	postMethod   = "POST"
	putMethod    = "PUT"

	// appleRegXMLString is the XML string for registering an iOS device
	// Replace {{Tags}} and {{DeviceID}} with the correct values
	appleRegXMLString string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
  <content type="application/xml">
    <AppleRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
      <Tags>{{Tags}}</Tags>
      <DeviceToken>{{DeviceID}}</DeviceToken>
    </AppleRegistrationDescription>
  </content>
</entry>`

	// appleTemplateRegXMLString is the XML string for registering an iOS device
	// Replace {{Tags}}, {{DeviceID}} and {{Template}} with the correct values
	appleTemplateRegXMLString string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
  <content type="application/xml">
    <AppleTemplateRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
      <Tags>{{Tags}}</Tags>
      <DeviceToken>{{DeviceID}}</DeviceToken>
			<BodyTemplate><![CDATA[{{Template}}]]></BodyTemplate>
    </AppleTemplateRegistrationDescription>
  </content>
</entry>`

	// gcmRegXMLString is the XML string for registering an iOS device
	// Replace {{Tags}} and {{DeviceID}} with the correct values
	gcmRegXMLString string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
  <content type="application/xml">
    <GcmRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
      <Tags>{{Tags}}</Tags>
      <GcmRegistrationId>{{DeviceID}}</GcmRegistrationId>
    </GcmRegistrationDescription>
  </content>
</entry>`

	// gcmRegTemplateXMLString is the XML string for registering an Android device
	// Replace {{Tags}}, {{DeviceID}} and {{Template}} with the correct values
	gcmTemplateRegXMLString string = `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom">
  <content type="application/xml">
    <GcmTemplateRegistrationDescription xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect">
      <Tags>{{Tags}}</Tags>
      <GcmRegistrationId>{{DeviceID}}</GcmRegistrationId>
			<BodyTemplate><![CDATA[{{Template}}]]></BodyTemplate>
    </GcmTemplateRegistrationDescription>
  </content>
</entry>`
)
