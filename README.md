# Azure Notification Hubs for Go(lang)

This library provides a Microsoft Azure Notification Hubs Client for backend applications.
It is packaged as a Go module to and is tested with Go 1.12+.

Originally a fork from [Gozure](https://github.com/onefootball/gozure) with patches
from [Martin Etnestad](https://github.com/gnawybol) @ [vippsas](https://github.com/vippsas/gozure).

Now maintained and packaged by [Daresay AB](https://daresay.co), [@daresaydigital](https://github.com/daresaydigital).

Basically a wrapper for this [Rest API](https://docs.microsoft.com/en-us/rest/api/notificationhubs/rest-api-methods)

[![Build Status](https://travis-ci.org/daresaydigital/azure-notificationhubs-go.svg?branch=master)](https://travis-ci.org/daresaydigital/azure-notificationhubs-go)

## Installing

Using go get

```sh
go get github.com/daresaydigital/azure-notificationhubs-go
```

## Registering device

```go
package main

import (
  "context"
  "github.com/daresaydigital/azure-notificationhubs-go"
)

func main() {
  var (
    hub      = notificationhubs.NewNotificationHub("YOUR_DefaultFullSharedAccessConnectionString", "YOUR_HubPath")
    template = `{
    "aps":{
      "alert":{
        "title":"$(title)",
        "body":"$(body)",
      },
      "badge":"#(badge)",
      "topic":"co.daresay.app",
      "content-available": 1
    },
    "name1":"$(value1)",
    "name2":"$(value2)"
  }`
  )

  template = strings.ReplaceAll(template, "\n", "")
  template = strings.ReplaceAll(template, "\t", "")

  reg := notificationhubs.NewTemplateRegistration(
    "ABC123",                       // The token from Apple or Google
    nil,                            // Expiration time, probably endless
    "ZXCVQWE",                      // Registration id, if you want to update an existing registration
    "tag1,tag2",                    // Tags that matches this device
    notificationhubs.ApplePlatform, // or GcmPlatform for Android
    template                        // The template. Use "$(name)" for strings and "#(name)" for numbers
  )

  // or hub.NewRegistration( ... ) without template

  hub.RegisterWithTemplate(context.TODO(), reg)
  // or if no template:
  hub.Register(context.TODO(), reg)
}
```

## Sending notification

```go
package main

import (
  "fmt"
  "github.com/daresaydigital/azure-notificationhubs-go"
)

func main() {
  var (
    hub     = notificationhubs.NewNotificationHub("YOUR_DefaultFullSharedAccessConnectionString", "YOUR_HubPath")
    payload = []byte(`{"title": "Hello Hub!"}`)
    n       = notificationhubs.NewNotification(notificationhubs.Template, payload)
  )

  // Broadcast push
  b, err := hub.Send(n, nil)
  if err != nil {
    panic(err)
  }

  fmt.Println("Message successfully created:", string(b))

  // Tag category push
  b, err = hub.Send(n, "tag1 || tag2")
  if err != nil {
    panic(err)
  }

  fmt.Println("Message successfully created:", string(b))
}
```

## Tag expressions

Read more about how to segment notification receivers in [the official documentation](https://docs.microsoft.com/en-us/azure/notification-hubs/notification-hubs-tags-segment-push-message).

### Example expressions

Example devices:

```json
"devices": {
  "A": {
    "tags": [
      "tag1",
      "tag2"
    ]
  },
  "B": {
    "tags": [
      "tag2",
      "tag3"
    ]
  },
  "C": {
    "tags": [
      "tag1",
      "tag2",
      "tag3"
    ]
  },
}
```

- Send to devices that has `tag1` or `tag2`. Example devices A, B and C.

  ```go
  hub.Send(notification, "tag1 || tag2")
  ```

- Send to devices that has `tag1` and `tag2`. Device A and C.

  ```go
  hub.Send(notification, "tag1 && tag2")
  ```

- Send to devices that has `tag1` and `tag2` but not `tag3`. Device A.

  ```go
  hub.Send(notification, "tag1 && tag2 && !tag3")
  ```

- Send to devices that has not `tag1`. Device B.

  ```go
  hub.Send(notification, "!tag1")
  ```

## Changelog

### v0.1.2

- Bugfix for reading the messsage id on standard hubs. Headers are always lowercase.

### v0.1.1

- Bugfix for when device registration were responding an unexpected response.

### v0.1.0

- Support for templated notifications
- Support for notification telemetry in higher tiers

### v0.0.2

- Big rewrite
- Added get registrations
- Travis CI
- Renamed the entities to use the same nomenclature as Azure
- Using fixtures for tests
- Support tag expressions

### v0.0.1

First release by Daresay. Restructured the code and renamed the API according to
Go standards.

## License

See the [LICENSE](LICENSE.txt) file for license rights and limitations (MIT).
