// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lovi-cloud/go-dorado-sdk/examples/lib"
)

func main() {
	client, err := lib.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	sps, err := client.LocalDevice.GetStoragePools(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range sps {
		fmt.Printf("%+v\n", v)
	}
}
