About
-----

The library provides a Microsoft Azure Notification Hub Client for backend applications.

Installing
----------

Using go get

```
$ go get github.com/onefootball/gozure/notihub
```

The package will be available under the following path:

```
$GOPATH/src/github.com/onefootball/gozure/notihub
```

Usage
-----

```go
package main

import (
    "fmt"
    "github.com/onefootball/gozure/notihub"
)

func main() {
    payload := []byte(`{"title": "Hello Hub!"}`)

    n, err := notihub.NewNotification(notihub.Template, payload)
    if err != nil {
        panic(err)
    }

    hub := notihub.NewNotificationHub("YOUR_DefaultFullSharedAccessSignature", "YOUR_HubPath")

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