# go-dorado-sdk: Golang SDK for Huawei OceanStor Dorado

go-dorado-sdk is Golang SDK for Huawei OceanStor Dorado.
It provides operationally from Golang.

## Support Products

- OceanStor Dorado V3

## Usage

- Enable REST API
- Create system-user that call REST API

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

## Reference documents

- [Developer Documents by Huawei](https://support.huawei.com/enterprise/en/centralized-storage/oceanstor-dorado3000-v3-pid-23786734?category=developer-documents)
    - require credential