# Azure Notification Hubs for Go(lang)

This library provides a Microsoft Azure Notification Hubs Client for backend applications.
It is packaged as a Go module to and is tested with Go 1.11+.

Originally a fork from [Gozure](https://github.com/onefootball/gozure) with patches
from [Martin Etnestad](https://github.com/gnawybol).

Now maintained and packaged by [Daresay AB](https://daresay.co), [@daresaydigital](https://github.com/daresaydigital).

## Installing

Using go get

```sh
$ go get github.com/daresaydigital/azure-notificationhubs-go
```

## Usage

```go
package notificationhubs

import (
    "fmt"
    "github.com/daresaydigital/azure-notificationhubs-go"
)

func main() {
    payload := []byte(`{"title": "Hello Hub!"}`)

    n, err := notificationhubs.NewNotification(notificationhubs.Template, payload)
    if err != nil {
        panic(err)
    }

    hub := notificationhubs.NewNotificationHub("YOUR_DefaultFullSharedAccessConnectionString", "YOUR_HubPath")

    // broadcast push
    b, err := hub.Send(n, nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("message successfully created:", string(b))

    // tag category push
    b, err = hub.Send(n, []string{"tag1", "tag2"})
    if err != nil {
        panic(err)
    }

    fmt.Println("message successfully created:", string(b))
}
```

## Changelog

### v0.0.1
First release by Daresay. Restructured the code and renamed the API according to
Go standards.

## License

See the [LICENSE](LICENSE.txt) file for license rights and limitations (MIT).
