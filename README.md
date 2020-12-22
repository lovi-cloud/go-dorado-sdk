# go-dorado-sdk

## Usage

```go
package main

import (
    "github.com/lovi-cloud/go-dorado-sdk/dorado"
)

func main() {
    username := "username"
	password := "password"
	localIps := []string{"https://192.0.2.100:8088", "https://192.0.2.101:8088"}
	remoteIps := []string{"https://192.0.2.200:8088", "https://192.0.2.201:8088"}
	portgroupName := "Port_Group"

	client, err := dorado.NewClient(localIps, remoteIps, username, password, portgroupName, nil)

    // etc...
}
```