About
-----

The library provides a Microsoft Azure Notification Hubs Client for backend applications.
It is packaged as a Go module

Installing
----------

Using go get

```
$ go get github.com/daresaydigital/azure-notificationhubs-go
```

The package will be available under the following path:

```
$GOPATH/src/github.com/daresaydigital/azure-notificationhubs-go
```

Usage
-----

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

License
-------
See the [LICENSE](LICENSE.txt) file for license rights and limitations (MIT).
